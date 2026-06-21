package repository

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/ledger/internal/domain"
	"github.com/jackc/pgx/v5"
)

type LedgerRepository interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
	CreateLedgerEntry(ctx context.Context, tx pgx.Tx, entry *domain.Ledger) error
	FindByID(ctx context.Context, id string) (*domain.Ledger, error)
	FindByReference(ctx context.Context, reference string) (*domain.Ledger, error)
	FindAll(ctx context.Context) ([]domain.Ledger, error)
	FindByEntryType(ctx context.Context, entryType string) ([]domain.Ledger, error)
	SumByAccountID(ctx context.Context, accID string) (int64, error)
	SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error
}
