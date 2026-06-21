package eventhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/domain"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type TransferConsumer struct {
	repo      repository.TransferRepository
	idemStore idempotency.Store
	publisher *producer.Publisher
	logger    *logrus.Logger
}

func NewTransferConsumer(repo repository.TransferRepository, idemStore idempotency.Store, publisher *producer.Publisher, logger *logrus.Logger) *TransferConsumer {
	return &TransferConsumer{
		repo:      repo,
		idemStore: idemStore,
		publisher: publisher,
		logger:    logger,
	}
}

func (h *TransferConsumer) OnTransferInitiated(ctx context.Context, msg []byte) error {
	h.logger.Info("TRANSFER INITIATED RECEIVED")
	var event events.TransferInitiatedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, event.Reference)
		if err != nil {
			return err
		}

		if processed {
			return nil // idempotent
		}

		err = h.repo.SaveOutboxEvent(ctx, tx, events.TopicDebitRequested, event.Reference,
			events.DebitRequestedEvent{
				BaseEvent: events.NewBaseEvent(ctx, "transfer-service"),
				Reference: event.Reference,
				FromAccID: event.FromAccID,
				FromAccNo: event.FromAccNo,
				ToAccID:   event.ToAccID,
				ToAccNo:   event.ToAccNo,
				Amount:    event.Amount,
			},
		)
		if err != nil {
			return err
		}

		return h.idemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), event.Reference, events.TopicTransferInitiated)
	})
}

func (h *TransferConsumer) OnDebitCompleted(ctx context.Context, msg []byte) error {
	h.logger.Info("TRANSFER DEBIT COMPLETED")
	var event events.DebitCompletedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}
	h.logger.Infof("EVENT ENTERED: DebitCompleted reference=%s amount=%d", event.Reference, event.Amount)

	eventID := fmt.Sprintf("%s:%s", event.Reference, events.TopicDebitCompleted)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		h.logger.Info("CHECKING IDEMPOTENCY")

		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			return err
		}

		h.logger.Infof("processed=%v", processed)

		if processed {
			return nil // idempotent
		}
		h.logger.Infof("IDEMPOTENCY CHECK PASSED reference=%s", eventID)

		h.logger.Info("UPDATING STEP")
		if err := h.repo.UpdateStep(ctx, tx, event.Reference, domain.StepDebited); err != nil {
			return err
		}

		h.logger.Info("UPDATING BALANCE")
		if err := h.repo.UpdateAfterDebit(ctx, tx, event.Reference, event.FromBalanceAfter); err != nil {
			return err
		}
		h.logger.Infof("DB UPDATED AFTER DEBIT reference=%s balance=%d", event.Reference, event.FromBalanceAfter)

		h.logger.Info("FETCHING TRANSFER")
		transfer, err := h.repo.FindByReference(ctx, tx, event.Reference)
		if err != nil {
			return err
		}

		h.logger.Info("SAVING CREDIT REQUEST")
		err = h.repo.SaveOutboxEvent(ctx, tx, events.TopicCreditRequested, event.Reference, events.CreditRequestedEvent{
			BaseEvent: events.NewBaseEvent(ctx, "transfer-service"),
			Reference: event.Reference,
			FromAccID: event.FromAccID,
			FromAccNo: transfer.FromAccNo,
			ToAccID:   transfer.ToAccID,
			ToAccNo:   transfer.ToAccNo,
			Amount:    event.Amount,
		})
		if err != nil {
			h.logger.WithError(err).Error("FAILED TO SAVE OUTBOX")
			return err
		}

		h.logger.Info("MARKING EVENT PROCESSED")
		return h.idemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), eventID, events.TopicDebitCompleted)
	})

}

func (h *TransferConsumer) OnDebitFailed(ctx context.Context, msg []byte) error {
	h.logger.Info("TRANSFER DEBIT FAILED")
	var event events.DebitFailedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	eventID := fmt.Sprintf("%s:%s", event.Reference, events.TopicDebitFailed)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			return err
		}
		if processed {
			return nil // idempotent
		}
		// update transfer status to failed
		if err := h.repo.UpdateStatus(ctx, tx, event.Reference, string(domain.TransferStatusFailed)); err != nil {
			return err
		}

		return h.idemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), eventID, events.TopicDebitFailed)
	})
}

func (h *TransferConsumer) OnCreditCompleted(ctx context.Context, msg []byte) error {
	h.logger.Info("TRANSFER CREDIT COMPLETED")
	var event events.CreditCompletedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}
	h.logger.Infof("EVENT ENTERED: CreditCompleted reference=%s amount=%d", event.Reference, event.Amount)

	eventID := fmt.Sprintf("%s:%s", event.Reference, events.TopicCreditCompleted)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			return err
		}
		if processed {
			return nil // idempotent
		}
		h.logger.Infof("IDEMPOTENCY CHECK PASSED reference=%s", eventID)

		if err := h.repo.UpdateAfterCredit(ctx, tx, event.Reference, event.ToBalanceAfter); err != nil {
			return err
		}

		if err := h.repo.UpdateStatus(ctx, tx, event.Reference, string(domain.TransferStatusSuccess)); err != nil {
			return err
		}
		h.logger.Infof("DB UPDATED AFTER CREDIT reference=%s balance=%d", event.Reference, event.ToBalanceAfter)

		transfer, err := h.repo.FindByReference(ctx, tx, event.Reference)
		if err != nil {
			return err
		}
		h.logger.Infof("PUBLISH CREDIT REQUEST reference=%s", event.Reference)

		err = h.repo.SaveOutboxEvent(ctx, tx, events.TopicTransferCompleted, event.Reference, events.TransferCompletedEvent{
			BaseEvent:        events.NewBaseEvent(ctx, "transfer-service"),
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
			return err
		}

		return h.idemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), eventID, events.TopicCreditCompleted)
	})
}

func (h *TransferConsumer) OnCreditFailed(ctx context.Context, msg []byte) error {
	h.logger.Info("TRANSFER CREDIT FAILED")
	var event events.CreditFailedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	eventID := fmt.Sprintf("%s:%s", event.Reference, events.TopicCreditFailed)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			return err
		}
		if processed {
			return nil // idempotent
		}

		transfer, err := h.repo.FindByReference(ctx, tx, event.Reference)
		if err != nil {
			return err
		}

		failedEvent := events.ReverseDebitEvent{
			BaseEvent: events.NewBaseEvent(ctx, "transfer-service"),
			Reference: transfer.Reference,
			AccountID: transfer.FromAccID,
			AccountNo: transfer.FromAccNo,
			Amount:    transfer.Amount,
		}

		// reverse debit
		if err = h.repo.SaveOutboxEvent(ctx, tx, events.TopicCreditFailed, event.Reference, failedEvent); err != nil {
			return err
		}

		if err := h.repo.UpdateStatus(ctx, tx, event.Reference, string(domain.TransferStatusFailed)); err != nil {
			return err
		}

		return h.idemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), eventID, events.TopicCreditFailed)
	})
}

func (h *TransferConsumer) OnDebitReverseCompleted(ctx context.Context, msg []byte) error {
	h.logger.Info("TRANSFER DEBIT REVERSE COMPLETED")
	var event events.ReverseDebitEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	eventID := fmt.Sprintf("%s:%s", event.Reference, events.TopicDebitReverseCompleted)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			return err
		}
		if processed {
			return nil // idempotent
		}

		// mark transfer as reversed
		if err := h.repo.UpdateStatus(ctx, tx, event.Reference, string(domain.TransferStatusReversing)); err != nil {
			return err
		}

		if err := h.repo.MarkAsReversed(ctx, tx, event.Reference); err != nil {
			return err
		}

		return h.idemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), eventID, events.TopicDebitReverseCompleted)
	})
}
