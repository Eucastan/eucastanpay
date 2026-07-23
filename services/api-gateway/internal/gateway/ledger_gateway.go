package gateway

import (
	"context"
	ledgerpb "github.com/Eucastan/eucastanpay/common/proto/ledger"
	ledgerReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/ledger"
	ledgerResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/ledger"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/mapper"
)

type LedgerGateway struct {
	client ledgerpb.LedgerServiceClient
}

func NewLedgerGateway(client ledgerpb.LedgerServiceClient) *LedgerGateway {
	return &LedgerGateway{
		client: client,
	}
}

func (s *LedgerGateway) GetAllLedgers(ctx context.Context) ([]*ledgerResp.LedgerResponse, error) {
	grpcResp, err := s.client.GetAllLedgers(
		ctx,
		mapper.ToProtoListLedgers(),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToListLedgerResponse(grpcResp)
	return resp, nil
}

func (s *LedgerGateway) GetLedgerUserID(ctx context.Context, userID string) (*ledgerResp.LedgerResponse, error) {
	grpcResp, err := s.client.GetLedgerByUserId(
		ctx,
		mapper.ToProtoLedgerByUserID(userID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToLedgerResponse(grpcResp)
	return resp, nil
}

func (s *LedgerGateway) GetLedger(ctx context.Context, ledgerID string) (*ledgerResp.LedgerResponse, error) {
	grpcResp, err := s.client.GetLedger(
		ctx,
		mapper.ToProtoLedger(ledgerID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToLedgerResponse(grpcResp)
	return resp, nil
}

func (s *LedgerGateway) GetLedgerBalance(ctx context.Context, accID string) (*ledgerResp.AccountBalanceResponse, error) {
	grpcResp, err := s.client.GetLedgerBalance(
		ctx,
		mapper.ToProtoLedgerBalance(accID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToLedgerBalanceResponse(accID, grpcResp)
	return resp, nil
}

func (s *LedgerGateway) GetLedgerByAccountID(ctx context.Context, accID string) (*ledgerResp.LedgerResponse, error) {
	grpcResp, err := s.client.GetLedgerByAccountId(
		ctx,
		mapper.ToProtoLedgerByAccountID(accID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToLedgerResponse(grpcResp)
	return resp, nil
}

func (s *LedgerGateway) GetLedgersByEntryType(ctx context.Context, input *ledgerReq.EntryTypeRequest) ([]*ledgerResp.LedgerResponse, error) {
	grpcResp, err := s.client.GetLedgersByEntryType(
		ctx,
		mapper.ToProtoLedgerEntryType(input),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToListLedgerEntryTypeResponse(grpcResp)
	return resp, nil
}

func (s *LedgerGateway) GetLedgerReconciliation(ctx context.Context, accID string) (*ledgerResp.ReconciliationResultResponse, error) {
	grpcResp, err := s.client.ReconcileAccount(
		ctx,
		mapper.ToProtoReconciliationByAccountID(accID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToReconciliationResponse(grpcResp)
	return resp, nil
}
