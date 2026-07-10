package usecase

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/audit/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository/postgres"
)

type AuditUseCase interface {
	Search(ctx context.Context, f postgres.Filter) ([]response.AuditReadResponse, error)
	GetAuditReadByID(ctx context.Context, id string) (*response.AuditReadResponse, error)
}
