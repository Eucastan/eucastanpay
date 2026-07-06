package eventhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/account/internal/domain"
	"github.com/Eucastan/eucastanpay/services/account/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/account/internal/repository"
	"github.com/Eucastan/eucastanpay/services/account/internal/usecase"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type AccountConsumer struct {
	Repo       repository.AccountRepository
	AccUseCase usecase.AccountUseCase
	IdemStore  idempotency.Store
	Publisher  *producer.Publisher
	telemetry  *telemetry.Telemetry
	logger     *logrus.Logger
}

func NewAccountConsumer(
	repo repository.AccountRepository,
	accUseCase usecase.AccountUseCase,
	idemStore idempotency.Store,
	publisher *producer.Publisher,
	telemetry *telemetry.Telemetry,
	logger *logrus.Logger,
) *AccountConsumer {
	return &AccountConsumer{
		Repo:       repo,
		AccUseCase: accUseCase,
		IdemStore:  idemStore,
		Publisher:  publisher,
		telemetry:  telemetry,
		logger:     logger,
	}
}

func (h *AccountConsumer) OnUserRegistration(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "AccountConsumer.OnUserRegistration")
	defer span.End()

	var event events.UserRegisteredEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		span.RecordError(err)
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": event.UserID,
		"email":   event.Email,
	}).Info("user registration event")

	return h.Repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.IdemStore.IsEventProcessedTx(ctx, tx, event.UserID)
		if err != nil {
			span.RecordError(err)
			return err
		}
		if processed {
			return nil
		}

		createReq := &request.CreateAccountRequest{
			AccountType: string(domain.SavingsAccount),
			Currency:    "NGN",
		}

		acc, err := h.AccUseCase.CreateAccountTX(ctx, tx, event.UserID, event.Email, createReq)
		if err != nil {
			span.RecordError(err)
			failedEvents := events.CreateAccFailedEvent{
				EventMetadata: events.NewChildEvent(event.EventMetadata),
				AccountID:     "",
				Reason:        err.Error(),
			}

			if err := h.Repo.SaveOutboxEvent(ctx, tx, events.TopicCreateAccFailed, event.UserID, failedEvents); err != nil {
				span.RecordError(err)
				return err
			}

			h.logger.WithFields(logrus.Fields{
				"correlation_id": failedEvents.CorrelationID,
				"service":        failedEvents.CausationID,
				"reason":         failedEvents.Reason,
			})

			return err
		}

		err = h.Repo.SaveOutboxEvent(ctx, tx, events.TopicAccountCreated, event.UserID, events.AccountCreatedEvent{
			EventMetadata: events.NewChildEvent(event.EventMetadata),
			AccountID:     acc.ID,
			UserID:        event.UserID,
			Email:         event.Email,
			AccountNo:     acc.AccountNo,
			AccountType:   acc.AccountType,
			Currency:      acc.Currency,
			Timestamp:     acc.CreatedAt.Unix(),
		})
		if err != nil {
			span.RecordError(err)
			return err
		}

		h.logger.WithFields(logrus.Fields{
			"user_id":      acc.UserID,
			"account_id":   acc.ID,
			"account_type": acc.AccountType,
		}).Info("user account creation successful")

		return h.IdemStore.MarkEventProcessedTx(
			ctx,
			tx,
			uuid.NewString(),
			event.UserID,
			events.TopicUserRegistered,
		)
	})
}

func (h *AccountConsumer) OnDebitRequested(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "AccountConsumer.OnDebitRequested")
	defer span.End()

	var event events.DebitRequestedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		span.RecordError(err)
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"operation":  "OnDebiteRequest",
		"account_no": event.FromAccNo,
		"amount":     event.Amount,
	}).Info("DEBIT REQUESTED EVENT:", event)

	input := &request.DebitRequest{
		AccountNo: event.FromAccNo,
		Amount:    event.Amount,
	}

	eventID := fmt.Sprintf("%s:%s", event.Reference, events.TopicDebitRequested)
	return h.Repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.IdemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			return err
		}
		if processed {
			return nil
		}

		h.logger.WithFields(logrus.Fields{
			"operation": "OnDebitRequested",
			"reference": event.Reference,
		}).Info("Processing debit request")

		if err := h.AccUseCase.Debit(ctx, tx, event.FromAccID, input); err != nil {
			failEvent := events.DebitFailedEvent{
				EventMetadata: events.NewChildEvent(event.EventMetadata),
				Reference:     event.Reference,
				Reason:        err.Error(),
			}

			if err := h.Repo.SaveOutboxEvent(ctx, tx, events.TopicDebitFailed, event.Reference, failEvent); err != nil {
				span.RecordError(err)
				return err
			}

			return err
		}
		h.logger.Info("Debit request processed successfully")

		acc, err := h.Repo.FindByIDTX(ctx, tx, event.FromAccID, event.FromAccNo)
		if err != nil {
			span.RecordError(err)
			return err
		}

		h.logger.Info("Saving debited data to outbox to publish later for debit completed")
		err = h.Repo.SaveOutboxEvent(ctx, tx, events.TopicDebitCompleted,
			event.Reference,
			events.DebitCompletedEvent{
				EventMetadata:    events.NewChildEvent(event.EventMetadata),
				FromAccID:        event.FromAccID,
				Reference:        event.Reference,
				Amount:           event.Amount,
				FromBalanceAfter: acc.Balance,
				Timestamp:        time.Now().Unix(),
			},
		)
		if err != nil {
			span.RecordError(err)
			return err
		}
		h.logger.Info("Debited data saved and ready to be published")

		return h.IdemStore.MarkEventProcessedTx(
			ctx,
			tx,
			uuid.NewString(),
			eventID,
			events.TopicDebitRequested,
		)

	})
}

func (h *AccountConsumer) OnCreditRequested(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "AccountConsumer.OnCreditRequested")
	defer span.End()

	var event events.CreditRequestedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		span.RecordError(err)
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"operation":  "OnCreditRequested",
		"account_no": event.ToAccNo,
		"amount":     event.Amount,
	}).Info("CREDIT REQUESTED EVENT:", event)

	input := &request.CreditRequest{
		AccountNo: event.ToAccNo,
		Amount:    event.Amount,
	}

	h.logger.Info("Before entering transaction")

	eventID := fmt.Sprintf("%s:%s", event.Reference, events.TopicCreditRequested)
	return h.Repo.WithTx(ctx, func(tx pgx.Tx) error {
		h.logger.Info("Entered transaction and about to check idempotency")
		processed, err := h.IdemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("Failed to check idempotency")
			return err
		}
		if processed {
			return nil
		}
		h.logger.Info("Idempotency check passed")

		h.logger.Info("Before ConfirmAccountNo")

		confirmAcc, err := h.Repo.ConfirmAccountNo(ctx, tx, event.ToAccNo)
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("ConfirmAccountNo failed")
			return err
		}

		h.logger.Info("After ConfirmAccountNo")

		h.logger.WithFields(logrus.Fields{
			"operation":            "OnCreditRequested",
			"reference":            event.Reference,
			"confirmed_account_id": confirmAcc.ID,
		}).Info("Processing credit request")

		h.logger.Info("Before Credit")

		if err := h.AccUseCase.Credit(ctx, tx, confirmAcc.ID, input); err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("Credit failed")
			failEvent := events.CreditFailedEvent{
				EventMetadata: events.NewChildEvent(event.EventMetadata),
				Reference:     event.Reference,
				Reason:        err.Error(),
			}
			if err := h.Repo.SaveOutboxEvent(ctx, tx, events.TopicCreditFailed, event.Reference, failEvent); err != nil {
				span.RecordError(err)
				return err
			}
			return err
		}
		h.logger.Info("After Credit")

		h.logger.Info("Credit request processed successfully")

		acc, err := h.Repo.FindByIDTX(ctx, tx, confirmAcc.ID, event.ToAccNo)
		if err != nil {
			span.RecordError(err)
			return err
		}

		h.logger.Info("Saving credited data to outbox...")

		err = h.Repo.SaveOutboxEvent(ctx, tx, events.TopicCreditCompleted, event.Reference, events.CreditCompletedEvent{
			EventMetadata:  events.NewChildEvent(event.EventMetadata),
			ToAccID:        confirmAcc.ID,
			Reference:      event.Reference,
			Amount:         event.Amount,
			ToBalanceAfter: acc.Balance,
			Timestamp:      time.Now().Unix(),
		})
		if err != nil {
			span.RecordError(err)
			return err
		}

		h.logger.Info("Saving credited data to outbox completed successfully")

		return h.IdemStore.MarkEventProcessedTx(
			ctx, tx,
			uuid.NewString(),
			eventID,
			events.TopicCreditRequested,
		)
	})
}
