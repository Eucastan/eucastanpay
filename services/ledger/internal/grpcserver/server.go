package grpcserver

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/proto/ledger"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LedgerServiceServer struct {
	ledger.UnimplementedLedgerServiceServer
	Ledger usecase.LedgerUseCase
}

func NewLedgerServiceServer(ledger usecase.LedgerUseCase) *LedgerServiceServer {
	return &LedgerServiceServer{Ledger: ledger}
}

func (s *LedgerServiceServer) ReconcileAccount(ctx context.Context, req *ledger.ReconcileAccountRequest) (*ledger.ReconcileResponse, error) {
	result, err := s.Ledger.ReconcileAccount(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	return &ledger.ReconcileResponse{
		Status:     result.Status,
		Difference: result.Difference,
		Message:    result.Reason,
	}, nil
}

func (s *LedgerServiceServer) GetAllLedgers(
	ctx context.Context,
	req *ledger.ListLedgersRequest,
) (*ledger.ListLedgersResponse, error) {
	ledgers, err := s.Ledger.GetAllLedgers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get all ledgers")
	}

	resp := make([]*ledger.Ledger, 0, len(ledgers))
	for _, v := range ledgers {
		resp = append(resp, &ledger.Ledger{
			LedgerId:     v.ID,
			AccountId:    v.AccountID,
			Amount:       v.Amount,
			EntryType:    v.EntryType,
			Reference:    v.Reference,
			BalanceAfter: v.BalanceAfter,
			Description:  v.Description,
			CreatedAt:    timestamppb.New(v.CreatedAt),
		})
	}

	return &ledger.ListLedgersResponse{
		Ledgers: resp,
	}, nil
}

func (s *LedgerServiceServer) GetLedgersByEntryType(
	ctx context.Context,
	req *ledger.ListLedgersEntryTypeRequest,
) (*ledger.ListLedgersEntryTypeResponse, error) {

	entryType := request.EntryTypeRequest{EntryType: req.EntryType}

	ledgers, err := s.Ledger.GetTransactionByEntryType(ctx, &entryType)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get all ledgers")
	}

	resp := make([]*ledger.Ledger, 0, len(ledgers))
	for _, v := range ledgers {
		resp = append(resp, &ledger.Ledger{
			LedgerId:     v.ID,
			AccountId:    v.AccountID,
			Amount:       v.Amount,
			EntryType:    v.EntryType,
			Reference:    v.Reference,
			BalanceAfter: v.BalanceAfter,
			Description:  v.Description,
			CreatedAt:    timestamppb.New(v.CreatedAt),
		})
	}

	return &ledger.ListLedgersEntryTypeResponse{
		Ledgers: resp,
	}, nil
}

func (s *LedgerServiceServer) GetLedgerByID(ctx context.Context, req *ledger.LedgerRequest) (*ledger.LedgerResponse, error) {
	ldg, err := s.Ledger.GetTransactionEntry(ctx, req.AccountId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch ledger entry")
	}

	return &ledger.LedgerResponse{
		LedgerId:     ldg.ID,
		AccountId:    ldg.AccountID,
		Amount:       ldg.Amount,
		EntryType:    ldg.EntryType,
		Reference:    ldg.Reference,
		BalanceAfter: ldg.BalanceAfter,
		Description:  ldg.Description,
		CreatedAt:    timestamppb.New(ldg.CreatedAt),
	}, nil
}
