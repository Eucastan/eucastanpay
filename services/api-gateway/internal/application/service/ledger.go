package service

import (
	"context"

	ledgerReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/ledger"
	ledgerResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/ledger"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/gateway"
)

type LedgerApplication struct {
	gateway *gateway.LedgerGateway
}

func NewLedgerApplication(gateway *gateway.LedgerGateway) *LedgerApplication {
	return &LedgerApplication{
		gateway: gateway,
	}
}

func (s *LedgerApplication) GetAllLedgers(ctx context.Context) ([]*ledgerResp.LedgerResponse, error) {
	return s.gateway.GetAllLedgers(ctx)
}

func (s *LedgerApplication) GetLedger(ctx context.Context, ledgerID string) (*ledgerResp.LedgerResponse, error) {
	return s.gateway.GetLedger(ctx, ledgerID)
}

func (s *LedgerApplication) GetLedgerUserID(ctx context.Context, userID string) (*ledgerResp.LedgerResponse, error) {
	return s.gateway.GetLedgerUserID(ctx, userID)
}

func (s *LedgerApplication) GetLedgerBalance(ctx context.Context, accID string) (*ledgerResp.AccountBalanceResponse, error) {
	return s.gateway.GetLedgerBalance(ctx, accID)
}

func (s *LedgerApplication) GetLedgerByAccountID(ctx context.Context, accID string) (*ledgerResp.LedgerResponse, error) {
	return s.gateway.GetLedgerByAccountID(ctx, accID)
}

func (s *LedgerApplication) GetLedgersByEntryType(ctx context.Context, input *ledgerReq.EntryTypeRequest) ([]*ledgerResp.LedgerResponse, error) {
	return s.gateway.GetLedgersByEntryType(ctx, input)
}

func (s *LedgerApplication) GetLedgerReconciliation(ctx context.Context, accID string) (*ledgerResp.ReconciliationResultResponse, error) {
	return s.gateway.GetLedgerReconciliation(ctx, accID)
}
