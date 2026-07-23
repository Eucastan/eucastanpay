package service

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/audit/internal/dto/response"
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

func (u *AuditUseCase) Search(ctx context.Context, f postgres.Filter) ([]response.AuditReadResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "AuditUseCase.Search")
	defer span.End()

	data, err := u.Repo.Search(ctx, f)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	resp := make([]response.AuditReadResponse, 0, len(data))
	for _, read := range data {
		resp = append(resp, *response.ToAuditReadResponse(&read))
	}

	return resp, nil
}

func (u *AuditUseCase) GetAuditReadByID(ctx context.Context, id string) (*response.AuditReadResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "AuditUseCase.GetAuditReadByID")
	defer span.End()

	auditRead, err := u.Repo.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return response.ToAuditReadResponse(auditRead), nil
}

func (u *AuditUseCase) GetAllAuditReads(ctx context.Context) ([]response.AuditReadResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "AuditUseCase.GetAllAuditReads")
	defer span.End()

	reads, err := u.Repo.FindAllAuditRead(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	resp := make([]response.AuditReadResponse, 0, len(reads))
	for _, read := range reads {
		resp = append(resp, *response.ToAuditReadResponse(&read))
	}

	return resp, nil
}
