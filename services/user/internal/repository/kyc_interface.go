package repository

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/user/internal/domain"
	"github.com/jackc/pgx/v5"
)

type KYCRepository interface {
	WithTX(ctx context.Context, fn func(tx pgx.Tx) error) error
	Create(ctx context.Context, kyc *domain.KYC) error
	FindByID(ctx context.Context, id string) (*domain.KYC, error)
	FindByUserID(ctx context.Context, userID string) (*domain.KYC, error)
	Update(ctx context.Context, kyc *domain.KYC) error
	SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error
}
