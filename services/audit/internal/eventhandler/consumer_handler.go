package eventhandler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/services/audit/internal/domain"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type AuditConsumer struct {
	Repo      repository.AuditRepository
	IdemStore idempotency.Store
	logger    *logrus.Logger
}

func NewAuditConsumer(repo repository.AuditRepository, idemStore idempotency.Store, log *logrus.Logger) *AuditConsumer {
	return &AuditConsumer{
		Repo:      repo,
		IdemStore: idemStore,
		logger:    log,
	}
}

func (c *AuditConsumer) Handler(topic string) func(ctx context.Context, msg []byte) error {
	return func(ctx context.Context, msg []byte) error {
		return c.handle(ctx, topic, msg)
	}
}

func (c *AuditConsumer) handle(ctx context.Context, topic string, msg []byte) error {
	c.logger.Info("Entered Audit Consumer Handler")

	payload := map[string]interface{}{}
	if err := json.Unmarshal(msg, &payload); err != nil {
		return err
	}

	log := &domain.AuditLog{
		ID:            uuid.NewString(),
		EventType:     topic,
		CorrelationID: getString(payload, "correlation_id"),
		Reference:     getString(payload, "reference"),
		Payload:       payload,
		CreatedAt:     time.Now(),
	}

	c.logger.Infof("payload=%+v", payload)

	eventID := getString(payload, "correlation_id")
	return c.Repo.WithTX(ctx, func(tx pgx.Tx) error {
		processed, err := c.IdemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			return err
		}
		if processed {
			return nil
		}

		if err := c.Repo.Insert(ctx, tx, log); err != nil {
			return err
		}

		if read := transformToRead(topic, payload); read != nil {
			c.logger.Infof(
				"transformed reference=%s correlation=%s",
				read.Reference,
				read.CorrelationID,
			)

			if err := c.Repo.InsertRead(ctx, tx, read); err != nil {
				return err
			}

			c.logger.Infof(
				"audit_read inserted reference=%s correlation=%s amount=%d",
				read.Reference,
				read.CorrelationID,
				read.Amount,
			)
		}

		return c.IdemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), eventID, topic)
	})
}
