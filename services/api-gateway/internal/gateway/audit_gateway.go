package gateway

import (
	"context"

	auditpb "github.com/Eucastan/eucastanpay/common/proto/audit"
	auditReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/audit"
	auditResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/audit"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/mapper"
)

type AuditGateway struct {
	client auditpb.AuditServiceClient
}

func NewAuditGateway(client auditpb.AuditServiceClient) *AuditGateway {
	return &AuditGateway{
		client: client,
	}
}

func (s *AuditGateway) GetAuditByID(ctx context.Context, auditID string) (*auditResp.AuditReadResponse, error) {
	grpcResp, err := s.client.GetAuditByID(
		ctx,
		mapper.ToProtoGetAuditIDRequest(auditID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToAuditReadResponse(grpcResp)
	return resp, nil
}

func (s *AuditGateway) GetAllAudits(ctx context.Context) (*auditResp.ReadResponse, error) {
	grpcResp, err := s.client.GetAllAudits(
		ctx,
		mapper.ToProtoListAudits(),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToListAuditReadResponse(grpcResp)
	return resp, nil
}

func (s *AuditGateway) SearchAuditLogs(ctx context.Context, input *auditReq.Filter) (*auditResp.ReadResponse, error) {

	grpcResp, err := s.client.SearchAudit(
		ctx,
		mapper.ToProtoSearchAuditLogs(input),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToSearchAuditResponse(grpcResp)
	return resp, nil
}
