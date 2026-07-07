package eventhandler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/audit/internal/domain"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type AuditConsumer struct {
	Repo      repository.AuditRepository
	IdemStore idempotency.Store
	telemetry *telemetry.Telemetry
	logger    *logrus.Logger
}

func NewAuditConsumer(repo repository.AuditRepository, idemStore idempotency.Store, telemetry *telemetry.Telemetry, log *logrus.Logger) *AuditConsumer {
	return &AuditConsumer{
		Repo:      repo,
		IdemStore: idemStore,
		telemetry: telemetry,
		logger:    log,
	}
}

func (c *AuditConsumer) Handler(topic string) func(ctx context.Context, msg []byte) error {
	return func(ctx context.Context, msg []byte) error {
		ctx, span := c.telemetry.Start(ctx, "AuditConsumer.Handler")
		defer span.End()

		return c.handle(ctx, topic, msg)
	}
}

func (c *AuditConsumer) handle(ctx context.Context, topic string, msg []byte) error {
	ctx, span := c.telemetry.Start(ctx, "AuditConsumer.handle")
	defer span.End()

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
			span.RecordError(err)
			return err
		}
		if processed {
			return nil
		}

		if err := c.Repo.Insert(ctx, tx, log); err != nil {
			span.RecordError(err)
			return err
		}

		if read := transformToRead(topic, payload); read != nil {
			c.logger.Infof(
				"transformed reference=%s correlation=%s",
				read.Reference,
				read.CorrelationID,
			)

			if err := c.Repo.InsertRead(ctx, tx, read); err != nil {
				span.RecordError(err)
				return err
			}

			c.logger.WithFields(logrus.Fields{
				"reference":      read.Reference,
				"correlation_id": read.CorrelationID,
				"amount":         read.Amount,
			}).Info("audit_read data inserted")
		}

		return c.IdemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), eventID, topic)
	})
}
