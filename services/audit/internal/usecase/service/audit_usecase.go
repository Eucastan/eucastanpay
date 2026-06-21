package service

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/audit/internal/domain"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository/postgres"
)

type AuditUseCase struct {
	Repo repository.AuditRepository
}

func NewAuditUseCase(repo repository.AuditRepository) *AuditUseCase {
	return &AuditUseCase{Repo: repo}
}

func (u *AuditUseCase) Search(ctx context.Context, f postgres.Filter) ([]domain.AuditRead, error) {
	return u.Repo.Search(ctx, f)
}

func (u *AuditUseCase) GetAuditReadByID(ctx context.Context, id string) (*domain.AuditRead, error) {
	auditLog, err := u.Repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return auditLog, nil
}
