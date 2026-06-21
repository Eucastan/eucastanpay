package usecase

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/audit/internal/domain"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository/postgres"
)

type AuditUseCase interface {
	Search(ctx context.Context, f postgres.Filter) ([]domain.AuditRead, error)
	GetAuditReadByID(ctx context.Context, id string) (*domain.AuditRead, error)
}
