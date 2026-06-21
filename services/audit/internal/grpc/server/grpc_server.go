package server

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/proto/audit"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/audit/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuditServiceServer struct {
	audit.UnimplementedAuditServiceServer
	Audit usecase.AuditUseCase
}

func NewAuditServiceServer(audit usecase.AuditUseCase) *AuditServiceServer {
	return &AuditServiceServer{
		Audit: audit,
	}
}

func (s *AuditServiceServer) SearchAudit(ctx context.Context, req *audit.SearchRequest) (*audit.SearchResponse, error) {
	filter := postgres.Filter{
		CorrelationID: req.CorrelationId,
		Reference:     req.Reference,
		EventType:     req.EventType,
		MinAmount:     req.MinAmount,
		MaxAmount:     req.MaxAmount,
		Limit:         int(req.Limit),
		Offset:        int(req.Offset),
	}

	auditLogs, err := s.Audit.Search(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, "audit search failed")
	}

	resp := make([]*audit.AuditEntry, 0, len(auditLogs))
	for _, v := range auditLogs {
		resp = append(resp, &audit.AuditEntry{
			Id:            v.ID,
			EventType:     v.EventType,
			Service:       v.Service,
			CorrelationId: v.CorrelationID,
			Reference:     v.Reference,
			AccountId:     v.AccountID,
			UserId:        v.UserID,
			Amount:        v.Amount,
			Status:        v.Status,
			CreatedAt:     v.CreatedAt.Unix(),
		})
	}

	return &audit.SearchResponse{
		Entries:    resp,
		TotalCount: int32(len(resp)),
	}, nil
}

func (s *AuditServiceServer) GetAuditByID(ctx context.Context, req *audit.GetByIDRequest) (*audit.AuditEntry, error) {
	auditRead, err := s.Audit.GetAuditReadByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get audit read")
	}

	return &audit.AuditEntry{
		Id:            auditRead.ID,
		EventType:     auditRead.EventType,
		Service:       auditRead.Service,
		CorrelationId: auditRead.CorrelationID,
		Reference:     auditRead.Reference,
		AccountId:     auditRead.AccountID,
		UserId:        auditRead.UserID,
		Amount:        auditRead.Amount,
		Status:        auditRead.Status,
		CreatedAt:     auditRead.CreatedAt.Unix(),
	}, nil
}
