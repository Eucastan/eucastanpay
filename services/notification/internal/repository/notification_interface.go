package repository

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/notification/internal/domain"
	"github.com/jackc/pgx/v5"
)

type NotificationRepository interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
	CreateTx(ctx context.Context, tx pgx.Tx, n *domain.Notification) error
	Create(ctx context.Context, n *domain.Notification) error
	UpdateStatus(ctx context.Context, id string, status string) error
	FindByUserID(ctx context.Context, userID string, limit int) ([]domain.Notification, error)
	FindByReference(ctx context.Context, reference string) ([]domain.Notification, error)
	SaveTemplate(ctx context.Context, template *domain.NotificationTemplate) error
	FindTemplateByName(ctx context.Context, name string) (*domain.NotificationTemplate, error)
}
