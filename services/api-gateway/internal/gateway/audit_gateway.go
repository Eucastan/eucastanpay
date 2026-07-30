package gateway

import (
	"context"

	auditpb "github.com/Eucastan/eucastanpay/common/proto/audit"
	auditReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/audit"
	auditResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/audit"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/mapper"
	"github.com/sirupsen/logrus"
)

type AuditGateway struct {
	client auditpb.AuditServiceClient
	logger *logrus.Logger
}

func NewAuditGateway(client auditpb.AuditServiceClient, logger *logrus.Logger) *AuditGateway {
	return &AuditGateway{
		client: client,
		logger: logger,
	}
}

func (g *AuditGateway) GetAuditByID(ctx context.Context, auditID string) (*auditResp.AuditReadResponse, error) {
	g.logger.Infof("Gateway audit id = %q", auditID)

	req := mapper.ToProtoGetAuditIDRequest(auditID)

	g.logger.Infof("Proto audit id = %q", req.AuditId)

	grpcResp, err := g.client.GetAuditByID(ctx, req)
	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	g.logger.WithFields(logrus.Fields{
		"audit_id": grpcResp.AuditId,
	}).Info(grpcResp)

	resp := mapper.ToAuditReadResponse(grpcResp)
	return resp, nil
}

func (g *AuditGateway) GetAllAudits(ctx context.Context) (*auditResp.ReadResponse, error) {
	grpcResp, err := g.client.GetAllAudits(
		ctx,
		mapper.ToProtoListAudits(),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	g.logger.WithFields(logrus.Fields{
		"count": len(grpcResp.Data),
	}).Info(grpcResp)

	resp := mapper.ToListAuditReadResponse(grpcResp)
	return resp, nil
}

func (g *AuditGateway) SearchAuditLogs(ctx context.Context, input *auditReq.Filter) (*auditResp.ReadResponse, error) {

	grpcResp, err := g.client.SearchAudit(
		ctx,
		mapper.ToProtoSearchAuditLogs(input),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	g.logger.WithFields(logrus.Fields{
		"count": len(grpcResp.Entries),
	}).Info(grpcResp)

	resp := mapper.ToSearchAuditResponse(grpcResp)
	return resp, nil
}
