package eventhandler

import (
	"context"
	"fmt"
	"time"

	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/domain"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type TransferConsumer struct {
	repo      repository.TransferRepository
	idemStore idempotency.Store
	telemetry *telemetry.Telemetry
	publisher *producer.Publisher
	logger    *logrus.Logger
}

func NewTransferConsumer(repo repository.TransferRepository, idemStore idempotency.Store, telemetry *telemetry.Telemetry, publisher *producer.Publisher, logger *logrus.Logger) *TransferConsumer {
	return &TransferConsumer{
		repo:      repo,
		idemStore: idemStore,
		telemetry: telemetry,
		publisher: publisher,
		logger:    logger,
	}
}

func (h *TransferConsumer) OnTransferInitiated(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "TransferConsumer.OnTransferInitiated")
	defer span.End()

	h.logger.Info("Transfer Initiated Received")

	event, err := kafka.Decode[events.TransferInitiatedEvent](msg)
	if err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"operation": "OnTransferInitiated",
		"reference": event.Reference,
	}).Info("Entering the transactional block")

	return h.emitDebitRequest(ctx, SagaRequest{
		ParentMetadata: event.EventMetadata,
		UserID:         event.UserID,
		Reference:      event.Reference,
		FromAccID:      event.FromAccID,
		FromAccNo:      event.FromAccNo,
		ToAccID:        event.ToAccID,
		ToAccNo:        event.ToAccNo,
		Amount:         event.Amount,
		ProcessedTopic: events.TopicTransferInitiated,
	})
}

func (h *TransferConsumer) OnDebitCompleted(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "TransferConsumer.OnDebitCompleted")
	defer span.End()

	event, err := kafka.Decode[events.DebitCompletedEvent](msg)
	if err != nil {
		span.RecordError(err)
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"operation": "OnDebitCompleted",
		"amount":    event.Amount,
	}).Infof("Beginning OnDebitCompleted process")

	eventID := fmt.Sprintf("%s:%s", event.Reference, events.TopicDebitCompleted)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			return err
		}

		if processed {
			return events.ErrProcessed // idempotent
		}
		h.logger.Infof("Idempotency Check Passed reference=%s", eventID)

		h.logger.Info("Before Update Step")
		if err := h.repo.UpdateStep(ctx, tx, event.Reference, string(domain.StepDebited)); err != nil {
			span.RecordError(err)
			return err
		}

		h.logger.Info("Before Update Balance after Debit")
		if err := h.repo.UpdateAfterDebit(ctx, tx, event.Reference, event.FromBalanceAfter); err != nil {
			span.RecordError(err)
			return err
		}

		transfer, err := h.repo.FindByReference(ctx, tx, event.Reference)
		if err != nil {
			span.RecordError(err)
			return err
		}

		h.logger.Info("Saving for Credit Request")
		err = h.repo.SaveOutboxEvent(ctx, tx, events.TopicCreditRequested, event.Reference,
			events.CreditRequestedEvent{
				EventMetadata: events.NewChildEvent(event.EventMetadata),
				UserID:        event.UserID,
				Reference:     event.Reference,
				FromAccID:     event.FromAccID,
				FromAccNo:     transfer.FromAccNo,
				ToAccID:       transfer.ToAccID,
				ToAccNo:       transfer.ToAccNo,
				Amount:        event.Amount,
			})

		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("Failed to saved credit request to outbox")
			return err
		}

		h.logger.Info("Marking event as Processed")
		return h.idemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), eventID, events.TopicDebitCompleted)
	})

}

func (h *TransferConsumer) OnDebitFailed(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "TransferConsumer.OnDebitFailed")
	defer span.End()

	h.logger.Info("Debit Failed Event")

	event, err := kafka.Decode[events.DebitFailedEvent](msg)
	if err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"operation": "OnDebitFailed",
		"reference": event.Reference,
	}).Info("Processing 'OnDebitFailed' event")

	eventID := fmt.Sprintf("%s:%s", event.Reference, events.TopicDebitFailed)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			return err
		}
		if processed {
			return events.ErrProcessed // idempotent
		}

		if err := h.repo.UpdateStatus(ctx, tx, event.Reference, string(domain.TransferStatusFailed)); err != nil {
			span.RecordError(err)
			return err
		}

		return h.idemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), eventID, events.TopicDebitFailed)
	})
}

func (h *TransferConsumer) OnCreditCompleted(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "TransferConsumer.OnCreditCompleted")
	defer span.End()

	h.logger.Info("Credit Completed Event")

	event, err := kafka.Decode[events.CreditCompletedEvent](msg)
	if err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"operation": "OnCreditCompleted",
		"reference": event.Reference,
		"amount":    event.Amount,
	}).Info("Entering 'CreditCompleted' event")

	eventID := fmt.Sprintf("%s:%s", event.Reference, events.TopicCreditCompleted)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			return err
		}
		if processed {
			return events.ErrProcessed // idempotent
		}

		h.logger.Info("Before Update Balance after Credit")
		if err := h.repo.UpdateAfterCredit(ctx, tx, event.Reference, event.ToBalanceAfter); err != nil {
			h.logger.WithError(err).Error("failed to update balance after credit")
			return err
		}

		h.logger.Info("Before Update Step")
		if err := h.repo.UpdateStep(ctx, tx, event.Reference, string(domain.StepCredited)); err != nil {
			h.logger.WithError(err).Error("failed to update step after credit")
			return err
		}

		h.logger.Info("Before UpdateStatus")
		if err := h.repo.UpdateStatus(ctx, tx, event.Reference, string(domain.TransferStatusSuccess)); err != nil {
			return err
		}

		h.logger.Info("Before FindByReference")
		transfer, err := h.repo.FindByReference(ctx, tx, event.Reference)
		if err != nil {
			return err
		}

		h.logger.Info("Saving to outbox for Transfer Completed event")
		err = h.repo.SaveOutboxEvent(ctx, tx, events.TopicTransferCompleted, event.Reference, events.TransferCompletedEvent{
			EventMetadata:    events.NewChildEvent(event.EventMetadata),
			TransferID:       transfer.ID,
			Reference:        event.Reference,
			UserID:           transfer.UserID,
			FromAccID:        transfer.FromAccID,
			ToAccID:          event.ToAccID,
			Amount:           event.Amount,
			FromBalanceAfter: transfer.FromBalanceAfter,
			ToBalanceAfter:   transfer.ToBalanceAfter,
			Timestamp:        time.Now().Unix(),
		})
		if err != nil {
			span.RecordError(err)
			return err
		}

		h.logger.Info("Successfully saved to outbox for transfer completed event")

		h.logger.Info("Marking event as Processed")
		return h.idemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), eventID, events.TopicCreditCompleted)
	})
}

func (h *TransferConsumer) OnCreditFailed(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "TransferConsumer.OnCreditFailed")
	defer span.End()

	h.logger.Info("Credit Failed Event")

	event, err := kafka.Decode[events.CreditFailedEvent](msg)
	if err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"operation": "OnCreditFailed",
		"reference": event.Reference,
	}).Info("Processing 'CreditFailed' event")

	eventID := fmt.Sprintf("%s:%s", event.Reference, events.TopicCreditFailed)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			return err
		}
		if processed {
			return events.ErrProcessed // idempotent
		}

		transfer, err := h.repo.FindByReference(ctx, tx, event.Reference)
		if err != nil {
			span.RecordError(err)
			return err
		}

		failedEvent := events.ReverseFailedTransferEvent{
			EventMetadata: events.NewChildEvent(event.EventMetadata),
			UserID:        transfer.UserID,
			Reference:     transfer.Reference,
			AccountID:     transfer.FromAccID,
			AccountNo:     transfer.FromAccNo,
			Amount:        transfer.Amount,
		}

		// reverse debit
		if err = h.repo.SaveOutboxEvent(ctx, tx, events.TopicReverseInitiated, event.Reference, failedEvent); err != nil {
			span.RecordError(err)
			return err
		}

		if err := h.repo.UpdateStatus(ctx, tx, event.Reference, string(domain.TransferStatusFailed)); err != nil {
			span.RecordError(err)
			return err
		}

		return h.idemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), eventID, events.TopicCreditFailed)
	})
}
