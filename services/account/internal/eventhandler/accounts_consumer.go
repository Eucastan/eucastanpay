package eventhandler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
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
	logger     *logrus.Logger
}

func NewAccountConsumer(repo repository.AccountRepository, accUseCase usecase.AccountUseCase, idemStore idempotency.Store, publisher *producer.Publisher, logger *logrus.Logger) *AccountConsumer {
	return &AccountConsumer{
		Repo:       repo,
		AccUseCase: accUseCase,
		IdemStore:  idemStore,
		Publisher:  publisher,
		logger:     logger,
	}
}

func (h *AccountConsumer) OnUserRegistration(ctx context.Context, msg []byte) error {
	var event events.UserRegisteredEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": event.UserID,
		"email":   event.Email,
	}).Info("user registration event")

	return h.Repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.IdemStore.IsEventProcessedTx(ctx, tx, event.UserID)
		if err != nil {
			return err
		}
		if processed {
			return nil
		}

		createReq := &request.CreateAccountRequest{
			AccountType: string(domain.SavingsAccount),
			Currency:    "NGN",
		}

		acc, err := h.AccUseCase.CreateAccountTX(ctx, tx, event.UserID, createReq)
		if err != nil {

			failedEvents := events.CreateAccFailedEvent{
				BaseEvent: events.NewBaseEvent(ctx, "account-service"),
				AccountID: "",
				Reason:    err.Error(),
			}
			if err := h.Repo.SaveOutboxEvent(ctx, tx, events.TopicCreateAccFailed, event.UserID, failedEvents); err != nil {
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
			BaseEvent:   events.NewBaseEvent(ctx, "account-service"),
			AccountID:   acc.ID,
			UserID:      acc.UserID,
			AccountNo:   acc.AccountNo,
			AccountType: acc.AccountType,
			Currency:    acc.Currency,
			Timestamp:   acc.CreatedAt.Unix(),
		})
		if err != nil {
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
	h.logger.Info("DEBIT REQUESTED MESSAGE:", string(msg))

	var event events.DebitRequestedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	h.logger.Info("DEBIT REQUESTED EVENT:", event)

	input := &request.DebitRequest{
		AccountNo: event.FromAccNo,
		Amount:    event.Amount,
	}

	return h.Repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.IdemStore.IsEventProcessedTx(ctx, tx, event.Reference)
		if err != nil {
			return err
		}
		if processed {
			return nil
		}

		if err := h.AccUseCase.Debit(ctx, tx, event.FromAccID, input); err != nil {
			failEvent := events.DebitFailedEvent{
				BaseEvent: events.NewBaseEvent(ctx, "account-service"),
				Reference: event.Reference,
				Reason:    err.Error(),
			}
			_ = h.Repo.SaveOutboxEvent(ctx, tx, events.TopicDebitFailed, event.Reference, failEvent)
			return err
		}

		acc, err := h.Repo.FindByIDTX(ctx, tx, event.FromAccID, event.FromAccNo)
		if err != nil {
			return err
		}

		err = h.Repo.SaveOutboxEvent(ctx, tx, events.TopicDebitCompleted,
			event.Reference,
			events.DebitCompletedEvent{
				BaseEvent:        events.NewBaseEvent(ctx, "account-service"),
				FromAccID:        event.FromAccID,
				Reference:        event.Reference,
				Amount:           event.Amount,
				FromBalanceAfter: acc.Balance,
				Timestamp:        time.Now().Unix(),
			},
		)
		if err != nil {
			return err
		}

		return h.IdemStore.MarkEventProcessedTx(
			ctx,
			tx,
			uuid.NewString(),
			event.Reference,
			events.TopicDebitRequested,
		)

	})
}

func (h *AccountConsumer) OnCreditRequested(ctx context.Context, msg []byte) error {
	h.logger.Info("CREDIT REQUESTED MESSAGE:", string(msg))

	var event events.CreditRequestedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	h.logger.Info("CREDIT REQUESTED EVENT:", event)

	input := &request.CreditRequest{
		AccountNo: event.ToAccNo,
		Amount:    event.Amount,
	}

	return h.Repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.IdemStore.IsEventProcessedTx(ctx, tx, event.Reference)
		if err != nil {
			return err
		}
		if processed {
			return nil
		}

		if err := h.AccUseCase.Credit(ctx, tx, event.ToAccID, input); err != nil {
			failEvent := events.CreditFailedEvent{
				BaseEvent: events.NewBaseEvent(ctx, "account-service"),
				Reference: event.Reference,
				Reason:    err.Error(),
			}
			_ = h.Repo.SaveOutboxEvent(ctx, tx, events.TopicCreditFailed, event.Reference, failEvent)
			return err
		}

		acc, err := h.Repo.FindByIDTX(ctx, tx, event.ToAccID, event.ToAccNo)
		if err != nil {
			return err
		}

		err = h.Repo.SaveOutboxEvent(ctx, tx, events.TopicCreditCompleted, event.Reference, events.CreditCompletedEvent{
			BaseEvent:      events.NewBaseEvent(ctx, "account-service"),
			ToAccID:        event.ToAccID,
			Reference:      event.Reference,
			Amount:         event.Amount,
			ToBalanceAfter: acc.Balance,
			Timestamp:      time.Now().Unix(),
		})
		if err != nil {
			return err
		}

		return h.IdemStore.MarkEventProcessedTx(
			ctx, tx,
			uuid.NewString(),
			event.Reference,
			events.TopicCreditRequested,
		)
	})
}
