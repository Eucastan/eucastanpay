package service

import (
	"context"

	auditReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/audit"
	auditResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/audit"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/gateway"
)

type AuditApplication struct {
	gateway *gateway.AuditGateway
}

func NewAuditApplication(gateway *gateway.AuditGateway) *AuditApplication {
	return &AuditApplication{
		gateway: gateway,
	}
}

func (s *AuditApplication) GetAuditByID(ctx context.Context, auditID string) (*auditResp.AuditReadResponse, error) {
	return s.gateway.GetAuditByID(
		ctx,
		auditID,
	)
}

func (s *AuditApplication) GetAllAudits(ctx context.Context) (*auditResp.ReadResponse, error) {
	return s.gateway.GetAllAudits(ctx)
}

func (s *AuditApplication) SearchAuditLogs(ctx context.Context, input *auditReq.Filter) (*auditResp.ReadResponse, error) {

	return s.gateway.SearchAuditLogs(
		ctx,
		input,
	)
}
