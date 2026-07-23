package grpcserver

import (
	"context"

	ledgerpb "github.com/Eucastan/eucastanpay/common/proto/ledger"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LedgerServiceServer struct {
	ledgerpb.UnimplementedLedgerServiceServer
	Ledger usecase.LedgerUseCase
}

func NewLedgerServiceServer(ledger usecase.LedgerUseCase) *LedgerServiceServer {
	return &LedgerServiceServer{Ledger: ledger}
}

func (s *LedgerServiceServer) ReconcileAccount(ctx context.Context, req *ledgerpb.ReconcileAccountRequest) (*ledgerpb.ReconcileResponse, error) {
	result, err := s.Ledger.ReconcileAccount(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	return &ledgerpb.ReconcileResponse{
		Status:         result.Status,
		AccountId:      result.AccountID,
		AccountBalance: result.AccountBalance,
		LedgerBalance:  result.LedgerBalance,
		Difference:     result.Difference,
		Message:        result.Reason,
		ReconciledAt:   timestamppb.New(result.ReconciledAt),
	}, nil
}

func (s *LedgerServiceServer) GetAllLedgers(
	ctx context.Context,
	req *ledgerpb.ListLedgersRequest,
) (*ledgerpb.ListLedgersResponse, error) {
	ledgers, err := s.Ledger.GetAllLedgers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get all ledgers")
	}

	resp := make([]*ledgerpb.Ledger, 0, len(ledgers))
	for _, v := range ledgers {
		resp = append(resp, &ledgerpb.Ledger{
			LedgerId:     v.ID,
			UserId:       v.UserID,
			AccountId:    v.AccountID,
			Amount:       v.Amount,
			EntryType:    v.EntryType,
			Reference:    v.Reference,
			BalanceAfter: v.BalanceAfter,
			Description:  v.Description,
			CreatedAt:    timestamppb.New(v.CreatedAt),
		})
	}

	return &ledgerpb.ListLedgersResponse{
		Ledgers: resp,
	}, nil
}

func (s *LedgerServiceServer) GetLedgersByEntryType(
	ctx context.Context,
	req *ledgerpb.ListLedgersEntryTypeRequest,
) (*ledgerpb.ListLedgersEntryTypeResponse, error) {

	entryType := request.EntryTypeRequest{EntryType: req.EntryType}

	ledgers, err := s.Ledger.GetTransactionByEntryType(ctx, &entryType)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get all ledgers")
	}

	resp := make([]*ledgerpb.Ledger, 0, len(ledgers))
	for _, v := range ledgers {
		resp = append(resp, &ledgerpb.Ledger{
			LedgerId:     v.ID,
			UserId:       v.UserID,
			AccountId:    v.AccountID,
			Amount:       v.Amount,
			EntryType:    v.EntryType,
			Reference:    v.Reference,
			BalanceAfter: v.BalanceAfter,
			Description:  v.Description,
			CreatedAt:    timestamppb.New(v.CreatedAt),
		})
	}

	return &ledgerpb.ListLedgersEntryTypeResponse{
		Ledgers: resp,
	}, nil
}

func (s *LedgerServiceServer) GetLedgerByID(ctx context.Context, req *ledgerpb.LedgerRequest) (*ledgerpb.LedgerResponse, error) {
	ldg, err := s.Ledger.GetLedger(ctx, req.LedgerId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch ledger entry")
	}

	return &ledgerpb.LedgerResponse{
		LedgerId:     ldg.ID,
		UserId:       ldg.UserID,
		AccountId:    ldg.AccountID,
		Amount:       ldg.Amount,
		EntryType:    ldg.EntryType,
		Reference:    ldg.Reference,
		BalanceAfter: ldg.BalanceAfter,
		Description:  ldg.Description,
		CreatedAt:    timestamppb.New(ldg.CreatedAt),
	}, nil
}

func (s *LedgerServiceServer) GetLedgerByUserID(ctx context.Context, req *ledgerpb.UserIdRequest) (*ledgerpb.LedgerResponse, error) {
	ldg, err := s.Ledger.GetLedger(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch ledger entry")
	}

	return &ledgerpb.LedgerResponse{
		LedgerId:     ldg.ID,
		UserId:       ldg.UserID,
		AccountId:    ldg.AccountID,
		Amount:       ldg.Amount,
		EntryType:    ldg.EntryType,
		Reference:    ldg.Reference,
		BalanceAfter: ldg.BalanceAfter,
		Description:  ldg.Description,
		CreatedAt:    timestamppb.New(ldg.CreatedAt),
	}, nil
}

func (s *LedgerServiceServer) GetLedgerByAccountID(ctx context.Context, req *ledgerpb.LedgerByAccountIdRequest) (*ledgerpb.LedgerResponse, error) {
	ldg, err := s.Ledger.GetTransactionEntry(ctx, req.AccountId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch ledger entry")
	}

	return &ledgerpb.LedgerResponse{
		LedgerId:     ldg.ID,
		UserId:       ldg.UserID,
		AccountId:    ldg.AccountID,
		Amount:       ldg.Amount,
		EntryType:    ldg.EntryType,
		Reference:    ldg.Reference,
		BalanceAfter: ldg.BalanceAfter,
		Description:  ldg.Description,
		CreatedAt:    timestamppb.New(ldg.CreatedAt),
	}, nil
}

func (s *LedgerServiceServer) GetBalance(ctx context.Context, req *ledgerpb.LedgerBalanceRequest) (*ledgerpb.LedgerBalanceResponse, error) {
	balance, err := s.Ledger.GetAccountBalance(ctx, req.AccountId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch ledger entry")
	}

	return &ledgerpb.LedgerBalanceResponse{
		Balance: balance,
	}, nil
}
