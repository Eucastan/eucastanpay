package repository

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/account/internal/domain"
	"github.com/jackc/pgx/v5"
)

type AccountRepository interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
	LockAccount(ctx context.Context, tx pgx.Tx, accID string, accNo int64) (*domain.Account, error)
	Create(ctx context.Context, tx pgx.Tx, acc *domain.Account) error
	Exists(ctx context.Context, tx pgx.Tx, userID, accType string) (bool, error)
	FindAll(ctx context.Context) ([]domain.Account, error)
	FindByID(ctx context.Context, accID string) (*domain.Account, error)
	ConfirmAccountNo(ctx context.Context, tx pgx.Tx, accNo int64) (*domain.Account, error)
	FindByIDTX(ctx context.Context, tx pgx.Tx, accID string, accNo int64) (*domain.Account, error)
	FindByUserID(ctx context.Context, userID string) (*domain.Account, error)
	FindByAccountIDAndUserID(ctx context.Context, accID, userID string) (*domain.Account, error)
	IsActive(ctx context.Context, tx pgx.Tx, accID, userID string) (bool, error)
	UpdateStatus(ctx context.Context, tx pgx.Tx, accID, status string) error
	UpdateBalance(ctx context.Context, tx pgx.Tx, id string, amount int64, isCredit bool) error
	SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error
}
