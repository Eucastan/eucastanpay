package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/notification/internal/domain"
	"github.com/Eucastan/eucastanpay/services/notification/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/notification/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/notification/internal/provider"
	"github.com/Eucastan/eucastanpay/services/notification/internal/repository"
)

type NotificationUseCase struct {
	repo      repository.NotificationRepository
	telemetry *telemetry.Telemetry
	providers *provider.EmailClient
	log       *logrus.Logger
}

func NewNotificationUseCase(repo repository.NotificationRepository, tm *telemetry.Telemetry, providers *provider.EmailClient, log *logrus.Logger) *NotificationUseCase {
	return &NotificationUseCase{
		repo:      repo,
		telemetry: tm,
		providers: providers,
		log:       log,
	}
}

// Send notification via preferred channel
func (u *NotificationUseCase) Send(ctx context.Context, email, userID string, req *request.NotificationRequest) error {
	ctx, span := u.telemetry.Start(ctx, "NotificationUseCase.Send")
	defer span.End()

	notification := &domain.Notification{
		ID:        uuid.NewString(),
		UserID:    userID,
		Title:     req.Title,
		Message:   req.Message,
		Channel:   req.Channel,
		Type:      req.Type,
		Priority:  req.Priority,
		Reference: req.Reference,
		Metadata:  req.Metadata,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	// Save to database first
	if err := u.repo.Create(ctx, notification); err != nil {
		return err
	}

	// Send via provider
	var sendErr error

	switch notification.Channel {
	case domain.ChannelEmail:
		sendErr = u.providers.SendEmail(notification)
	case domain.ChannelPush:
		sendErr = u.providers.SendPush(notification)
	case domain.ChannelInApp:
		sendErr = u.providers.SendInApp(notification)
	}

	if sendErr != nil {
		span.RecordError(sendErr)
		u.log.WithError(sendErr).WithField(
			"notification_id", notification.ID,
		).Error("Failed to send notification")
		notification.Status = "failed"
	} else {
		notification.Status = "sent"
		now := time.Now()
		notification.SentAt = &now
	}

	// Update final status
	u.repo.UpdateStatus(ctx, notification.ID, notification.Status)
	return sendErr
}

func (u *NotificationUseCase) GetByUserID(ctx context.Context, userID string) ([]response.NotificationResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "NotificationUseCase.GetByUserID")
	defer span.End()

	notifications, err := u.repo.FindByUserID(ctx, userID, 10)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	resp := response.ToNotificationResponse(notifications)
	return resp, nil
}

func (u *NotificationUseCase) GetByReference(ctx context.Context, ref string) ([]response.NotificationResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "NotificationUseCase.GetByReference")
	defer span.End()

	notifications, err := u.repo.FindByReference(ctx, ref)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	resp := response.ToNotificationResponse(notifications)
	return resp, nil
}

func (u *NotificationUseCase) SaveTemplate(ctx context.Context, tempt *request.NotificationTemplateRequest) error {
	ctx, span := u.telemetry.Start(ctx, "NotificationUseCase.SaveTemplate")
	defer span.End()

	template := &domain.NotificationTemplate{
		ID:      uuid.NewString(),
		Name:    tempt.Name,
		Subject: tempt.Subject,
		Body:    tempt.Body,
		Channel: tempt.Channel,
	}

	if err := u.repo.SaveTemplate(ctx, template); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (u *NotificationUseCase) GetByTemplateName(ctx context.Context, name string) (*response.NotificationTemplateResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "NotificationUseCase.GetByTemplateName")
	defer span.End()

	notifications, err := u.repo.FindTemplateByName(ctx, name)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	resp := response.ToNotificationTemplateResponse(notifications)
	return resp, nil
}
