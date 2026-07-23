package repository

import (
	"context"
	"time"

	"github.com/Eucastan/eucastanpay/services/transfer/internal/domain"
	"github.com/jackc/pgx/v5"
)

type TransferRepository interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
	Create(ctx context.Context, tx pgx.Tx, t *domain.Transfer) error
	FindAll(ctx context.Context) ([]domain.Transfer, error)
	FindByIdempotencyKey(ctx context.Context, idemKey string) (*domain.Transfer, error)
	FindByReference(ctx context.Context, tx pgx.Tx, ref string) (*domain.Transfer, error)
	FindByReferenceNoTx(ctx context.Context, reference string) (*domain.Transfer, error)
	FindByID(ctx context.Context, id string) (*domain.Transfer, error)
	FindByUserID(ctx context.Context, userID string) (*domain.Transfer, error)
	UpdateStatus(ctx context.Context, tx pgx.Tx, ref, status string) error
	UpdateAfterDebit(ctx context.Context, tx pgx.Tx, ref string, FromBalanceAfter int64) error
	UpdateAfterCredit(ctx context.Context, tx pgx.Tx, ref string, ToBalanceAfter int64) error
	UpdateStep(ctx context.Context, tx pgx.Tx, ref string, step string) error
	MarkAsReversed(ctx context.Context, tx pgx.Tx, ref string) error
	SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error
	IncrementRecoveryCount(ctx context.Context, tx pgx.Tx, reference string) error
	FindStuckTransfers(ctx context.Context, timeout time.Duration) ([]domain.Transfer, error)
}
