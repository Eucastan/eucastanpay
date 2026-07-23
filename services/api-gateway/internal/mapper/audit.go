package mapper

import (
	"encoding/json"

	auditpb "github.com/Eucastan/eucastanpay/common/proto/audit"
	auditReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/audit"

	auditResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/audit"
)

func ToProtoGetAuditIDRequest(auditID string) *auditpb.GetByIDRequest {
	return &auditpb.GetByIDRequest{
		AuditId: auditID,
	}
}

func ToProtoListAudits() *auditpb.AuditRequest {
	return &auditpb.AuditRequest{}
}

func ToProtoSearchAuditLogs(req *auditReq.Filter) *auditpb.SearchRequest {
	return &auditpb.SearchRequest{
		CorrelationId: req.CorrelationID,
		Reference:     req.Reference,
		EventType:     req.EventType,
		MinAmount:     req.MinAmount,
		MaxAmount:     req.MaxAmount,
		FromDate:      req.FromDate.Unix(),
		ToDate:        req.ToDate.Unix(),
		Limit:         int32(req.Limit),
		Offset:        int32(req.Offset),
	}
}

func ToAuditReadResponse(r *auditpb.AuditEntryResponse) *auditResp.AuditReadResponse {
	return &auditResp.AuditReadResponse{
		ID:            r.AuditId,
		EventType:     r.EventType,
		Service:       r.Service,
		CorrelationID: r.CorrelationId,
		CausationID:   r.CausationId,
		Reference:     r.Reference,
		AccountID:     r.AccountId,
		UserID:        r.UserId,
		Amount:        r.Amount,
		Status:        r.Status,
		Payload:       json.RawMessage(r.Payload),
		CreatedAt:     r.CreatedAt.AsTime(),
	}
}

func ToListAuditReadResponse(req *auditpb.GetAllAuditResponse) *auditResp.ReadResponse {
	data := make([]auditResp.AuditReadResponse, 0, len(req.Data))
	for _, r := range req.Data {
		data = append(data, auditResp.AuditReadResponse{
			ID:            r.AuditId,
			EventType:     r.EventType,
			Service:       r.Service,
			CorrelationID: r.CorrelationId,
			CausationID:   r.CausationId,
			Reference:     r.Reference,
			AccountID:     r.AccountId,
			UserID:        r.UserId,
			Amount:        r.Amount,
			Status:        r.Status,
			Payload:       json.RawMessage(r.Payload),
			CreatedAt:     r.CreatedAt.AsTime(),
		})
	}

	return &auditResp.ReadResponse{
		Data: data,
	}
}

func ToSearchAuditResponse(req *auditpb.SearchResponse) *auditResp.ReadResponse {
	data := make([]auditResp.AuditReadResponse, 0, len(req.Entries))
	for _, r := range req.Entries {
		data = append(data, auditResp.AuditReadResponse{
			ID:            r.AuditId,
			EventType:     r.EventType,
			Service:       r.Service,
			CorrelationID: r.CorrelationId,
			CausationID:   r.CausationId,
			Reference:     r.Reference,
			AccountID:     r.AccountId,
			UserID:        r.UserId,
			Amount:        r.Amount,
			Status:        r.Status,
			Payload:       json.RawMessage(r.Payload),
			CreatedAt:     r.CreatedAt.AsTime(),
		})
	}

	return &auditResp.ReadResponse{
		Data: data,
	}
}
