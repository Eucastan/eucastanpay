package server

import (
	"context"
	"encoding/json"
	"time"

	auditpb "github.com/Eucastan/eucastanpay/common/proto/audit"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/audit/internal/usecase"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuditServiceServer struct {
	auditpb.UnimplementedAuditServiceServer
	Audit  usecase.AuditUseCase
	logger *logrus.Logger
}

func NewAuditServiceServer(audit usecase.AuditUseCase, logger *logrus.Logger) *AuditServiceServer {
	return &AuditServiceServer{
		Audit:  audit,
		logger: logger,
	}
}

func (s *AuditServiceServer) SearchAudit(ctx context.Context, req *auditpb.SearchRequest) (*auditpb.SearchResponse, error) {
	s.logger.Info("AuditServiceServer.SearchAudit id reached >>")

	var fromDate *time.Time
	if req.FromDate != 0 {
		t := time.Unix(req.FromDate, 0)
		fromDate = &t
	}

	var toDate *time.Time
	if req.ToDate != 0 {
		t := time.Unix(req.ToDate, 0)
		toDate = &t
	}

	filter := postgres.Filter{
		CorrelationID: req.CorrelationId,
		Reference:     req.Reference,
		EventType:     req.EventType,
		MinAmount:     req.MinAmount,
		MaxAmount:     req.MaxAmount,
		FromDate:      fromDate,
		ToDate:        toDate,
		Limit:         int(req.Limit),
		Offset:        int(req.Offset),
	}

	auditLogs, err := s.Audit.Search(ctx, filter)
	if err != nil {
		s.logger.WithError(err).Error(err.Error())
		return nil, status.Error(codes.Internal, "audit search failed")
	}

	resp := make([]*auditpb.AuditEntryResponse, 0, len(auditLogs))
	for _, v := range auditLogs {
		var metadata map[string]any
		if err := json.Unmarshal(v.Payload, &metadata); err != nil {
			s.logger.WithError(err).Error(err.Error())
			return nil, status.Error(codes.Internal, "invalid payload")
		}

		payload, err := structpb.NewStruct(metadata)
		if err != nil {
			s.logger.WithError(err).Error(err.Error())
			return nil, status.Error(codes.Internal, "invalid metadata")
		}

		resp = append(resp, &auditpb.AuditEntryResponse{
			AuditId:       v.ID,
			EventType:     v.EventType,
			Service:       v.Service,
			CorrelationId: v.CorrelationID,
			CausationId:   v.CausationID,
			Reference:     v.Reference,
			AccountId:     v.AccountID,
			UserId:        v.UserID,
			Amount:        v.Amount,
			Status:        v.Status,
			Metadata:      payload,
			CreatedAt:     timestamppb.New(v.CreatedAt),
		})
	}

	s.logger.Info(resp)

	return &auditpb.SearchResponse{
		Entries:    resp,
		TotalCount: int32(len(resp)),
	}, nil
}

func (s *AuditServiceServer) GetAuditByID(ctx context.Context, req *auditpb.GetByIDRequest) (*auditpb.AuditEntryResponse, error) {
	s.logger.Info("AuditServiceServer.GetAuditByID is reached >>")

	auditRead, err := s.Audit.GetAuditReadByID(ctx, req.AuditId)
	if err != nil {
		s.logger.WithError(err).Error(err.Error())
		return nil, status.Error(codes.Internal, "failed to get audit read")
	}

	var metadata map[string]any
	if err := json.Unmarshal(auditRead.Payload, &metadata); err != nil {
		s.logger.WithError(err).Error(err.Error())
		return nil, status.Error(codes.Internal, "invalid payload")
	}

	payload, err := structpb.NewStruct(metadata)
	if err != nil {
		s.logger.WithError(err).Error(err.Error())
		return nil, status.Error(codes.Internal, "invalid metadata")
	}

	return &auditpb.AuditEntryResponse{
		AuditId:       auditRead.ID,
		EventType:     auditRead.EventType,
		Service:       auditRead.Service,
		CorrelationId: auditRead.CorrelationID,
		CausationId:   auditRead.CausationID,
		Reference:     auditRead.Reference,
		AccountId:     auditRead.AccountID,
		UserId:        auditRead.UserID,
		Amount:        auditRead.Amount,
		Status:        auditRead.Status,
		Metadata:      payload,
		CreatedAt:     timestamppb.New(auditRead.CreatedAt),
	}, nil
}

func (s *AuditServiceServer) GetAllAudits(ctx context.Context, req *auditpb.AuditRequest) (*auditpb.GetAllAuditResponse, error) {
	s.logger.Info("AuditServiceServer.GetAllAudits is reached >>")

	reads, err := s.Audit.GetAllAuditReads(ctx)
	if err != nil {
		s.logger.WithError(err).Error(err.Error())
		return nil, status.Error(codes.Internal, "failed to get audit read")
	}

	data := make([]*auditpb.AuditEntryResponse, 0, len(reads))
	for _, r := range reads {
		var metadata map[string]any
		if err := json.Unmarshal(r.Payload, &metadata); err != nil {
			s.logger.WithError(err).Error(err.Error())
			return nil, status.Error(codes.Internal, "invalid payload")
		}

		payload, err := structpb.NewStruct(metadata)
		if err != nil {
			s.logger.WithError(err).Error(err.Error())
			return nil, status.Error(codes.Internal, "invalid metadata")
		}

		data = append(data, &auditpb.AuditEntryResponse{
			AuditId:       r.ID,
			EventType:     r.EventType,
			Service:       r.Service,
			CorrelationId: r.CorrelationID,
			CausationId:   r.CausationID,
			Reference:     r.Reference,
			AccountId:     r.AccountID,
			UserId:        r.UserID,
			Amount:        r.Amount,
			Status:        r.Status,
			Metadata:      payload,
			CreatedAt:     timestamppb.New(r.CreatedAt),
		})
	}

	s.logger.WithFields(logrus.Fields{
		"count": len(data),
	}).Info(data)

	return &auditpb.GetAllAuditResponse{
		Data: data,
	}, nil
}
