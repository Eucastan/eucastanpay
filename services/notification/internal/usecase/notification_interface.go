package usecase

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/notification/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/notification/internal/dto/response"
)

type NotificationUseCase interface {
	Send(ctx context.Context, email, userID string, req *request.NotificationRequest) error
	GetByUserID(ctx context.Context, userID string) ([]response.NotificationResponse, error)
	GetByReference(ctx context.Context, ref string) ([]response.NotificationResponse, error)
	SaveTemplate(ctx context.Context, tempt *request.NotificationTemplateRequest) error
	GetByTemplateName(ctx context.Context, name string) (*response.NotificationTemplateResponse, error)
}
