package repository

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/admin/internal/domain"
	"github.com/jackc/pgx/v5"
)

type AdminRepository interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
	Create(ctx context.Context, admin *domain.Admin) error
	FindByEmail(ctx context.Context, email string) (*domain.Admin, error)
	FindByID(ctx context.Context, id string) (*domain.Admin, error)
	Update(ctx context.Context, admin *domain.Admin) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]domain.Admin, error)
	SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error
}
