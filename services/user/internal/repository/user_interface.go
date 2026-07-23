package repository

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/user/internal/domain"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	WithTX(ctx context.Context, fn func(tx pgx.Tx) error) error
	Create(ctx context.Context, user *domain.User) error
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, userID string) error
	SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error
}
