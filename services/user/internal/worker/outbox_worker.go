package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

const (
	maxRetries = 5
)

type OutboxWorker struct {
	db        *pgxpool.Pool
	publisher *producer.Publisher
	log       *logrus.Logger
	interval  time.Duration
}

type OutboxEvent struct {
	ID         string
	Topic      string
	Key        string
	Payload    []byte
	RetryCount int
}

func NewOutboxWorker(
	db *pgxpool.Pool,
	publisher *producer.Publisher,
	log *logrus.Logger,
	interval time.Duration,
) *OutboxWorker {

	if log == nil {
		log = logrus.New()
	}

	return &OutboxWorker{
		db:        db,
		publisher: publisher,
		log:       log,
		interval:  interval,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	w.log.Info("Outbox worker started")

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			w.log.Info("outbox worker stopped")
			return

		case <-ticker.C:
			w.processBatch(ctx)
		}
	}
}

func (w *OutboxWorker) fetchAndLock(ctx context.Context) ([]OutboxEvent, error) {
	query := `
	WITH cte AS (
    SELECT id
    FROM outbox
    WHERE published = false
      AND failed = false
      AND (locked_until IS NULL OR locked_until <= NOW())
    ORDER BY created_at
    LIMIT 50
    FOR UPDATE SKIP LOCKED
	)
	UPDATE outbox o
	SET locked_until = NOW() + INTERVAL '1 minute'
	FROM cte
	WHERE o.id = cte.id
	RETURNING o.id, o.topic, o.key, o.payload, o.retry_count;
	`

	rows, err := w.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []OutboxEvent

	for rows.Next() {

		var e OutboxEvent

		if err := rows.Scan(&e.ID, &e.Topic, &e.Key, &e.Payload,
			&e.RetryCount,
		); err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

func (w *OutboxWorker) processBatch(ctx context.Context) {

	events, err := w.fetchAndLock(ctx)
	if err != nil {
		w.log.WithError(err).Error("failed to fetch outbox events")
		return
	}

	w.log.Infof("found %d outbox events", len(events))

	for _, event := range events {
		w.log.Infof("publishing event topic=%s key=%s", event.Topic, event.Key)

		err := w.publisher.Publish(
			ctx,
			event.Topic,
			event.Key,
			json.RawMessage(event.Payload),
		)

		if err != nil {
			w.log.WithError(err).Error("publish failed")
			w.handleFailure(ctx, event.ID, event.RetryCount, err)

			continue
		}

		w.log.Infof("published topic=%s", event.Topic)
		w.markPublished(ctx, event.ID)

		w.log.WithFields(logrus.Fields{
			"event_id": event.ID,
			"topic":    event.Topic,
		}).Info("event published")
	}
}

func (w *OutboxWorker) markPublished(ctx context.Context, id string) {
	_, err := w.db.Exec(
		ctx,
		`
		UPDATE outbox
		SET published = true,
		    published_at = NOW(),
		    locked_until = NULL
		WHERE id = $1
		`,
		id,
	)

	if err != nil {
		w.log.WithError(err).Error("failed to mark event published")
	}
}

func (w *OutboxWorker) handleFailure(ctx context.Context, id string, retryCount int, publishErr error) {
	retryCount++

	if retryCount >= maxRetries {

		_, err := w.db.Exec(
			ctx,
			`
			UPDATE outbox
			SET failed = true,
			    retry_count = $2,
			    locked_until = NULL,
			    last_error = $3
			WHERE id = $1
			`,
			id,
			retryCount,
			publishErr.Error(),
		)

		if err != nil {
			w.log.WithError(err).Error("failed to move event to DLQ")
		}

		return
	}

	backoff := time.Second * time.Duration(1<<uint(retryCount))

	_, err := w.db.Exec(
		ctx,
		`
		UPDATE outbox
		SET retry_count = $2,
		    locked_until = NOW() + $3::interval,
		    last_error = $4
		WHERE id = $1
		`,
		id,
		retryCount,
		fmt.Sprintf("%f seconds", backoff.Seconds()),
		publishErr.Error(),
	)

	if err != nil {
		w.log.WithError(err).Error("failed to update retry state")
	}
}
