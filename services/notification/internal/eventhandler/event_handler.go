package eventhandler

import (
	"context"
	"fmt"
	"time"

	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/notification/internal/domain"
	"github.com/Eucastan/eucastanpay/services/notification/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type NotificationConsumer struct {
	repo      repository.NotificationRepository
	idemStore idempotency.Store
	telemetry *telemetry.Telemetry
	logger    *logrus.Logger
}

func NewNotificationConsumer(
	repo repository.NotificationRepository,
	idemStore idempotency.Store,
	telemetry *telemetry.Telemetry,
	logger *logrus.Logger,
) *NotificationConsumer {
	return &NotificationConsumer{
		repo:      repo,
		idemStore: idemStore,
		telemetry: telemetry,
		logger:    logger,
	}
}

func (h *NotificationConsumer) OnUserRegistered(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "NotificationConsumer.OnUserRegistered")
	defer span.End()

	e, err := kafka.Decode[events.UserRegisteredEvent](msg)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to unmarshal user registered event: %w", err)
	}

	eventID := fmt.Sprintf("%s:%s", e.UserID, events.TopicUserRegistered)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to check if event is processed")
			return err
		}
		if processed {
			h.logger.WithField("eventID", eventID).Info("event has already been processed")
			return events.ErrProcessed
		}

		err = h.repo.CreateTx(ctx, tx, &domain.Notification{
			ID:        uuid.NewString(),
			UserID:    e.UserID,
			Title:     "Welcome to EucastanPay!",
			Message:   fmt.Sprintf("Hi %s, welcome to EucastanPay. Your account has been created.", e.FirstName),
			Channel:   domain.ChannelEmail,
			Type:      domain.NotificationTypeAccount,
			CreatedAt: time.Now(),
		})
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to send notification")
			return err
		}

		return h.idemStore.MarkEventProcessedTx(
			ctx, tx,
			uuid.NewString(), eventID,
			events.TopicUserRegistered,
		)
	})
}

func (h *NotificationConsumer) OnKycCreated(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "NotificationConsumer.OnKycCreated")
	defer span.End()

	e, err := kafka.Decode[events.KYCCreatedEvent](msg)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).Error("failed to unmarshal KYC created event")
		return err
	}

	eventID := fmt.Sprintf("%s:%s", e.UserID, events.TopicUserKYCCreated)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to check if event is processed")
			return err
		}
		if processed {
			h.logger.WithField("eventID", eventID).Info("event has already been processed")
			return events.ErrProcessed
		}

		err = h.repo.CreateTx(ctx, tx, &domain.Notification{
			ID:        uuid.NewString(),
			UserID:    e.UserID,
			Title:     "KYC was created successfully",
			Message:   "Hi, welcome to EucastanPay. Your KYC has been created.",
			Channel:   domain.ChannelEmail,
			Type:      domain.NotificationTypeAccount,
			CreatedAt: time.Now(),
		})
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to send notification")
			return err
		}

		return h.idemStore.MarkEventProcessedTx(
			ctx, tx,
			uuid.NewString(), eventID,
			events.TopicUserKYCCreated,
		)
	})
}

func (h *NotificationConsumer) OnUserKycVerified(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "NotificationConsumer.UserKycVerified")
	defer span.End()

	e, err := kafka.Decode[events.UserKYCVerifiedEvent](msg)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).Error("failed to unmarshal KYC Verified event")
		return err
	}

	eventID := fmt.Sprintf("%s:%s", e.UserID, events.TopicUserKYCVerified)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to check if event is processed")
			return err
		}
		if processed {
			h.logger.WithField("eventID", eventID).Info("event has already been processed")
			return events.ErrProcessed
		}

		err = h.repo.CreateTx(ctx, tx, &domain.Notification{
			ID:        uuid.NewString(),
			UserID:    e.UserID,
			Title:     "KYC Verification",
			Message:   "Hi, welcome to EucastanPay. Your KYC has been verified.",
			Channel:   domain.ChannelEmail,
			Type:      domain.NotificationTypeAccount,
			CreatedAt: time.Now(),
		})
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to send notification")
			return err
		}

		return h.idemStore.MarkEventProcessedTx(
			ctx, tx,
			uuid.NewString(), eventID,
			events.TopicUserKYCVerified,
		)
	})
}

func (h *NotificationConsumer) OnAccountCreated(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "NotificationConsumer.OnAccountCreated")
	defer span.End()

	e, err := kafka.Decode[events.AccountCreatedEvent](msg)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).Error("failed to unmarshal account created event")
		return err
	}

	eventID := fmt.Sprintf("%s:%s", e.UserID, events.TopicAccountCreated)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to check if event is processed")
			return err
		}
		if processed {
			h.logger.WithField("eventID", eventID).Info("event has already been processed")
			return events.ErrProcessed
		}

		err = h.repo.CreateTx(ctx, tx, &domain.Notification{
			ID:        uuid.NewString(),
			UserID:    e.UserID,
			Title:     "Account was created successfully",
			Message:   "Hi, welcome to EucastanPay. Your account has been created.",
			Channel:   domain.ChannelEmail,
			Type:      domain.NotificationTypeAccount,
			CreatedAt: time.Now(),
		})
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to send notification")
			return err
		}

		return h.idemStore.MarkEventProcessedTx(
			ctx, tx,
			uuid.NewString(), eventID,
			events.TopicAccountCreated,
		)
	})
}

func (h *NotificationConsumer) OnAccountCreationFailed(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "NotificationConsumer.OnAccountCreationFailed")
	defer span.End()

	e, err := kafka.Decode[events.CreateAccFailedEvent](msg)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).Error("failed to unmarshal account creation failed event")
		return err
	}

	eventID := fmt.Sprintf("%s:%s", e.UserID, events.TopicCreateAccFailed)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to check if event is processed")
			return err
		}
		if processed {
			h.logger.WithField("eventID", eventID).Info("event has already been processed")
			return events.ErrProcessed
		}

		err = h.repo.CreateTx(ctx, tx, &domain.Notification{
			ID:        uuid.NewString(),
			UserID:    e.UserID,
			Title:     "Account Creation Failed",
			Message:   "Your account creation has failed.",
			Channel:   domain.ChannelEmail,
			Type:      domain.NotificationTypeAccount,
			CreatedAt: time.Now(),
		})
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to send notification")
			return err
		}

		return h.idemStore.MarkEventProcessedTx(
			ctx, tx,
			uuid.NewString(), eventID,
			events.TopicTransferFailed,
		)
	})
}

func (h *NotificationConsumer) OnAccountDeposit(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "NotificationConsumer.OnAccountDeposit")
	defer span.End()

	e, err := kafka.Decode[events.DepositAccountEvent](msg)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).Error("failed to unmarshal account deposit event")
		return err
	}

	eventID := fmt.Sprintf("%s:%s", e.Reference, events.TopicDepositAccount)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to check if event is processed")
			return err
		}
		if processed {
			h.logger.WithField("eventID", eventID).Info("event has already been processed")
			return events.ErrProcessed
		}

		err = h.repo.CreateTx(ctx, tx, &domain.Notification{
			ID:        uuid.NewString(),
			UserID:    e.UserID,
			Title:     "Account Deposit",
			Message:   fmt.Sprintf("%d was deposited into your account", e.Amount),
			Channel:   domain.ChannelEmail,
			Type:      domain.NotificationTypeTransaction,
			CreatedAt: time.Now(),
		})
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to send notification")
			return err
		}

		return h.idemStore.MarkEventProcessedTx(
			ctx, tx,
			uuid.NewString(), eventID,
			events.TopicDepositAccount,
		)
	})
}

func (h *NotificationConsumer) OnCashWithdraw(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "NotificationConsumer.OnCashWithdraw")
	defer span.End()

	e, err := kafka.Decode[events.WithdrawalEvent](msg)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).Error("failed to unmarshal cash withdraw event")
		return err
	}

	eventID := fmt.Sprintf("%s:%s", e.Reference, events.TopicWithdrawal)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to check if event is processed")
			return err
		}
		if processed {
			h.logger.WithField("eventID", eventID).Info("event has already been processed")
			return events.ErrProcessed
		}

		err = h.repo.CreateTx(ctx, tx, &domain.Notification{
			ID:        uuid.NewString(),
			UserID:    e.UserID,
			Title:     "Cash Withdraw Completed",
			Message:   fmt.Sprintf("%d was deducted from your account", e.Amount),
			Channel:   domain.ChannelEmail,
			Type:      domain.NotificationTypeTransaction,
			CreatedAt: time.Now(),
		})
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to send notification")
			return err
		}

		return h.idemStore.MarkEventProcessedTx(
			ctx, tx,
			uuid.NewString(), eventID,
			events.TopicWithdrawal,
		)
	})
}

func (h *NotificationConsumer) OnTransferCompleted(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "NotificationConsumer.OnTransferCompleted")
	defer span.End()

	e, err := kafka.Decode[events.TransferCompletedEvent](msg)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).Error("failed to unmarshal transfer completed event")
		return err
	}

	eventID := fmt.Sprintf("%s:%s", e.Reference, events.TopicTransferCompleted)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to check if event is processed")
			return err
		}
		if processed {
			h.logger.WithField("eventID", eventID).Info("event has already been processed")
			return events.ErrProcessed
		}

		err = h.repo.CreateTx(ctx, tx, &domain.Notification{
			ID:        uuid.NewString(),
			UserID:    e.UserID,
			Title:     "Transfer Completed",
			Message:   "Your transfer has been completed successfully.",
			Channel:   domain.ChannelEmail,
			Type:      domain.NotificationTypeTransaction,
			CreatedAt: time.Now(),
		})
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to send notification")
			return err
		}

		return h.idemStore.MarkEventProcessedTx(
			ctx, tx,
			uuid.NewString(), eventID,
			events.TopicTransferCompleted,
		)
	})
}

func (h *NotificationConsumer) OnTransferFailed(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "NotificationConsumer.OnTransferFailed")
	defer span.End()

	e, err := kafka.Decode[events.TransferFailedEvent](msg)
	if err != nil {
		span.RecordError(err)
		h.logger.WithError(err).Error("failed to unmarshal transfer failed event")
		return err
	}

	eventID := fmt.Sprintf("%s:%s", e.Reference, events.TopicTransferFailed)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to check if event is processed")
			return err
		}
		if processed {
			h.logger.WithField("eventID", eventID).Info("event has already been processed")
			return events.ErrProcessed
		}

		err = h.repo.CreateTx(ctx, tx, &domain.Notification{
			ID:        uuid.NewString(),
			UserID:    e.UserID,
			Title:     "Transfer Failed",
			Message:   "Your transfer has failed.",
			Channel:   domain.ChannelEmail,
			Type:      domain.NotificationTypeTransaction,
			CreatedAt: time.Now(),
		})
		if err != nil {
			span.RecordError(err)
			h.logger.WithError(err).Error("failed to send notification")
			return err
		}

		return h.idemStore.MarkEventProcessedTx(
			ctx, tx,
			uuid.NewString(), eventID,
			events.TopicTransferFailed,
		)
	})
}
