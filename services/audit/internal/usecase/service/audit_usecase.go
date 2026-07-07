package service

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/audit/internal/domain"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository/postgres"
)

type AuditUseCase struct {
	Repo      repository.AuditRepository
	telemetry *telemetry.Telemetry
}

func NewAuditUseCase(repo repository.AuditRepository, telemetry *telemetry.Telemetry) *AuditUseCase {
	return &AuditUseCase{
		Repo:      repo,
		telemetry: telemetry,
	}
}

func (u *AuditUseCase) Search(ctx context.Context, f postgres.Filter) ([]domain.AuditRead, error) {
	ctx, span := u.telemetry.Start(ctx, "AuditUseCase.Search")
	defer span.End()

	data, err := u.Repo.Search(ctx, f)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return data, nil
}

func (u *AuditUseCase) GetAuditReadByID(ctx context.Context, id string) (*domain.AuditRead, error) {
	ctx, span := u.telemetry.Start(ctx, "AuditUseCase.GetAuditReadByID")
	defer span.End()

	auditLog, err := u.Repo.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return auditLog, nil
}
