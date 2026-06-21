package request

import "time"

type LedgerRequest struct {
	Amount    int64  `json:"amount" binding:"required"`
	EntryType string `json:"entry_type" binding:"required"`
	Reference string `json:"reference" binding:"required"`
}

type EntryTypeRequest struct {
	EntryType string `json:"entry_type,omitempty"`
}

type AccountBalanceFromLedger struct {
	AccountID string
	Balance   int64
}

type ReconciliationResult struct {
	AccountID      string    `json:"account_id"`
	AccountBalance int64     `json:"account_balance"`
	LedgerBalance  int64     `json:"ledger_balance"`
	Difference     int64     `json:"difference"`
	Status         string    `json:"status"` // "ok" or "discrepancy"
	Reason         string    `json:"reason,omitempty"`
	ReconciledAt   time.Time `json:"reconciled_at"`
}
