package service

import (
	"context"

	auditReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/audit"
	auditResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/audit"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/gateway"
	"github.com/sirupsen/logrus"
)

type AuditApplication struct {
	gateway *gateway.AuditGateway
	logger  *logrus.Logger
}

func NewAuditApplication(gateway *gateway.AuditGateway, logger *logrus.Logger) *AuditApplication {
	return &AuditApplication{
		gateway: gateway,
		logger:  logger,
	}
}

func (s *AuditApplication) GetAuditByID(ctx context.Context, auditID string) (*auditResp.AuditReadResponse, error) {
	s.logger.Infof("Application layer: %q", auditID)

	g, err := s.gateway.GetAuditByID(ctx, auditID)
	if err != nil {
		s.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	return g, nil
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
