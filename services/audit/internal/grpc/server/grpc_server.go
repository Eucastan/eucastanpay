package server

import (
	"context"

	auditpb "github.com/Eucastan/eucastanpay/common/proto/audit"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/audit/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuditServiceServer struct {
	auditpb.UnimplementedAuditServiceServer
	Audit usecase.AuditUseCase
}

func NewAuditServiceServer(audit usecase.AuditUseCase) *AuditServiceServer {
	return &AuditServiceServer{
		Audit: audit,
	}
}

func (s *AuditServiceServer) SearchAudit(ctx context.Context, req *auditpb.SearchRequest) (*auditpb.SearchResponse, error) {
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

	resp := make([]*auditpb.AuditEntryResponse, 0, len(auditLogs))
	for _, v := range auditLogs {
		resp = append(resp, &auditpb.AuditEntryResponse{
			AuditId:       v.ID,
			EventType:     v.EventType,
			Service:       v.Service,
			CorrelationId: v.CorrelationID,
			Reference:     v.Reference,
			AccountId:     v.AccountID,
			UserId:        v.UserID,
			Amount:        v.Amount,
			Status:        v.Status,
			CreatedAt:     timestamppb.New(v.CreatedAt),
		})
	}

	return &auditpb.SearchResponse{
		Entries:    resp,
		TotalCount: int32(len(resp)),
	}, nil
}

func (s *AuditServiceServer) GetAuditByID(ctx context.Context, req *auditpb.GetByIDRequest) (*auditpb.AuditEntryResponse, error) {
	auditRead, err := s.Audit.GetAuditReadByID(ctx, req.AuditId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get audit read")
	}

	return &auditpb.AuditEntryResponse{
		AuditId:       auditRead.ID,
		EventType:     auditRead.EventType,
		Service:       auditRead.Service,
		CorrelationId: auditRead.CorrelationID,
		Reference:     auditRead.Reference,
		AccountId:     auditRead.AccountID,
		UserId:        auditRead.UserID,
		Amount:        auditRead.Amount,
		Status:        auditRead.Status,
		CreatedAt:     timestamppb.New(auditRead.CreatedAt),
	}, nil
}

func (s *AuditServiceServer) GetAllAudits(ctx context.Context, req *auditpb.AuditRequest) (*auditpb.GetAllAuditResponse, error) {
	reads, err := s.Audit.GetAllAuditReads(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get audit read")
	}

	data := make([]*auditpb.AuditEntryResponse, 0, len(reads))
	for _, r := range reads {
		data = append(data, &auditpb.AuditEntryResponse{
			AuditId:       r.ID,
			EventType:     r.EventType,
			Service:       r.Service,
			CorrelationId: r.CorrelationID,
			Reference:     r.Reference,
			AccountId:     r.AccountID,
			UserId:        r.UserID,
			Amount:        r.Amount,
			Status:        r.Status,
			CreatedAt:     timestamppb.New(r.CreatedAt),
		})
	}

	return &auditpb.GetAllAuditResponse{
		Data: data,
	}, nil
}
