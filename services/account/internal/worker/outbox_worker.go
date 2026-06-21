package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

const (
	maxRetries   = 5
	lockDuration = 60 * time.Second // Lock for 1 minute
	backoffBase  = 2 * time.Second
)

type OutboxWorker struct {
	db        *pgxpool.Pool
	publisher *producer.Publisher
	log       *logrus.Logger
}

func NewOutboxWorker(db *pgxpool.Pool, publisher *producer.Publisher, log *logrus.Logger) *OutboxWorker {
	if log == nil {
		log = logrus.New()
	}
	return &OutboxWorker{
		db:        db,
		publisher: publisher,
		log:       log,
	}
}

func StartOutboxWorker(ctx context.Context, db *pgxpool.Pool, publisher *producer.Publisher, log *logrus.Logger) {
	worker := NewOutboxWorker(db, publisher, log)

	query := `
        SELECT id, topic, key, payload, retry_count
        FROM outbox 
        WHERE published = false 
          AND locked_until <= NOW()
        ORDER BY created_at ASC 
        LIMIT 50
    `

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			worker.log.Info("Outbox worker shutting down...")
			return
		case <-ticker.C:
			worker.processBatch(ctx, query)
		}
	}
}

func (w *OutboxWorker) processBatch(ctx context.Context, query string) {
	rows, err := w.db.Query(ctx, query)
	if err != nil {
		w.log.WithError(err).Error("Failed to query outbox")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, topic, key string
		var payload []byte
		var retryCount int

		if err := rows.Scan(&id, &topic, &key, &payload, &retryCount); err != nil {
			continue
		}

		// Lock the record immediately
		lockErr := w.lockRecord(ctx, id)
		if lockErr != nil {
			continue // Already locked by another instance
		}

		// Publish to Kafka
		err = w.publisher.Publish(ctx, topic, key, json.RawMessage(payload))
		if err != nil {
			w.handleFailure(ctx, id, retryCount)
			continue
		}

		// Success → Mark as published
		_, _ = w.db.Exec(ctx, `UPDATE outbox SET published = true WHERE id = $1`, id)
		w.log.WithField("event_id", id).Info("Outbox event published successfully")
	}
}

func (w *OutboxWorker) lockRecord(ctx context.Context, id string) error {
	_, err := w.db.Exec(ctx, `
        UPDATE outbox 
        SET locked_until = NOW() + $1 
        WHERE id = $2 
          AND locked_until <= NOW()
    `, lockDuration, id)
	return err
}

func (w *OutboxWorker) handleFailure(ctx context.Context, id string, retryCount int) {
	newRetry := retryCount + 1

	if newRetry >= maxRetries {
		w.log.WithField("event_id", id).Warn("Moving event to DLQ after max retries")
		// Optional: Move to DLQ table or just leave it marked with high retry count
		_, _ = w.db.Exec(ctx, `
            UPDATE outbox 
            SET locked_until = NOW() + INTERVAL '1 hour', 
                retry_count = $2 
            WHERE id = $1`, id, newRetry)
		return
	}

	// Exponential backoff
	backoff := time.Duration(newRetry) * backoffBase

	_, _ = w.db.Exec(ctx, `
        UPDATE outbox 
        SET locked_until = NOW() + $2, 
            retry_count = $3 
        WHERE id = $1`, id, backoff, newRetry)
}
