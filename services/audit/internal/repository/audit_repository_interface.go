package repository

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/audit/internal/domain"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository/postgres"
	"github.com/jackc/pgx/v5"
)

type AuditRepository interface {
	WithTX(ctx context.Context, fn func(tx pgx.Tx) error) error
	Insert(ctx context.Context, tx pgx.Tx, log *domain.AuditLog) error
	InsertRead(ctx context.Context, tx pgx.Tx, read *domain.AuditRead) error
	Search(ctx context.Context, f postgres.Filter) ([]domain.AuditRead, error)
	FindByID(ctx context.Context, id string) (*domain.AuditRead, error)
}
