package mapper

import (
	ledgerpb "github.com/Eucastan/eucastanpay/common/proto/ledger"
	ledgerReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/ledger"

	ledgerResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/ledger"
)

func ToProtoLedgerEntryType(req *ledgerReq.EntryTypeRequest) *ledgerpb.ListLedgersEntryTypeRequest {
	return &ledgerpb.ListLedgersEntryTypeRequest{
		EntryType: req.EntryType,
	}
}

func ToProtoListLedgers() *ledgerpb.ListLedgersRequest {
	return &ledgerpb.ListLedgersRequest{}
}

func ToProtoLedgerByUserID(userID string) *ledgerpb.UserIdRequest {
	return &ledgerpb.UserIdRequest{
		UserId: userID,
	}
}

func ToProtoLedger(ledgerID string) *ledgerpb.LedgerRequest {
	return &ledgerpb.LedgerRequest{
		LedgerId: ledgerID,
	}
}

func ToProtoLedgerBalance(accID string) *ledgerpb.LedgerBalanceRequest {
	return &ledgerpb.LedgerBalanceRequest{
		AccountId: accID,
	}
}

func ToProtoLedgerByAccountID(accID string) *ledgerpb.LedgerByAccountIdRequest {
	return &ledgerpb.LedgerByAccountIdRequest{
		AccountId: accID,
	}
}

func ToProtoReconciliationByAccountID(accID string) *ledgerpb.ReconcileAccountRequest {
	return &ledgerpb.ReconcileAccountRequest{
		AccountId: accID,
	}
}

func ToLedgerBalanceResponse(accID string, req *ledgerpb.LedgerBalanceResponse) *ledgerResp.AccountBalanceResponse {
	return &ledgerResp.AccountBalanceResponse{
		AccountID: accID,
		Balance:   req.Balance,
	}
}

func ToLedgerResponse(req *ledgerpb.LedgerResponse) *ledgerResp.LedgerResponse {
	return &ledgerResp.LedgerResponse{
		ID:           req.LedgerId,
		UserID:       req.UserId,
		AccountID:    req.AccountId,
		Amount:       req.Amount,
		EntryType:    req.EntryType,
		Reference:    req.Reference,
		BalanceAfter: req.BalanceAfter,
		Description:  req.Description,
		CreatedAt:    req.CreatedAt.AsTime(),
		UpdatedAt:    req.UpdatedAt.AsTime(),
	}
}

func ToListLedgerEntryTypeResponse(req *ledgerpb.ListLedgersEntryTypeResponse) []*ledgerResp.LedgerResponse {
	data := make([]*ledgerResp.LedgerResponse, 0, len(req.Ledgers))
	for _, r := range req.Ledgers {
		data = append(data, &ledgerResp.LedgerResponse{
			ID:           r.LedgerId,
			UserID:       r.UserId,
			AccountID:    r.AccountId,
			Amount:       r.Amount,
			EntryType:    r.EntryType,
			Reference:    r.Reference,
			BalanceAfter: r.BalanceAfter,
			Description:  r.Description,
			CreatedAt:    r.CreatedAt.AsTime(),
			UpdatedAt:    r.UpdatedAt.AsTime(),
		})
	}

	return data
}

func ToReconciliationResponse(req *ledgerpb.ReconcileResponse) *ledgerResp.ReconciliationResultResponse {
	return &ledgerResp.ReconciliationResultResponse{
		Data: ledgerResp.ReconciliationResult{
			AccountID:      req.AccountId,
			AccountBalance: req.AccountBalance,
			LedgerBalance:  req.LedgerBalance,
			Difference:     req.Difference,
			Status:         req.Status,
			Reason:         req.Message,
			ReconciledAt:   req.ReconciledAt.AsTime(),
		},
	}
}

func ToListLedgerResponse(req *ledgerpb.ListLedgersResponse) []*ledgerResp.LedgerResponse {
	data := make([]*ledgerResp.LedgerResponse, 0, len(req.Ledgers))
	for _, r := range req.Ledgers {
		data = append(data, &ledgerResp.LedgerResponse{
			ID:           r.LedgerId,
			UserID:       r.UserId,
			AccountID:    r.AccountId,
			Amount:       r.Amount,
			EntryType:    r.EntryType,
			Reference:    r.Reference,
			BalanceAfter: r.BalanceAfter,
			Description:  r.Description,
			CreatedAt:    r.CreatedAt.AsTime(),
			UpdatedAt:    r.UpdatedAt.AsTime(),
		})
	}

	return data
}
