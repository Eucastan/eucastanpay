package ledger

import "time"

type ReconciliationResult struct {
	AccountID      string    `json:"account_id"`
	AccountBalance int64     `json:"account_balance"`
	LedgerBalance  int64     `json:"ledger_balance"`
	Difference     int64     `json:"difference"`
	Status         string    `json:"status"` // "ok" or "discrepancy"
	Reason         string    `json:"reason,omitempty"`
	ReconciledAt   time.Time `json:"reconciled_at"`
}

type ReconciliationResultResponse struct {
	Data ReconciliationResult `json:"data"`
}

type AccountBalanceResponse struct {
	AccountID string `json:"account_id"`
	Balance   int64  `json:"balance"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type LedgerResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	AccountID    string    `json:"account_id"`
	Amount       int64     `json:"amount"`
	EntryType    string    `json:"entry_type"`
	Reference    string    `json:"reference"`
	BalanceAfter int64     `json:"balance_after"`
	Description  string    `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
