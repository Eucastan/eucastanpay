package domain

import (
	"time"
)

type LedgerEntryType string

const (
	DebitEntry  LedgerEntryType = "debit"
	CreditEntry LedgerEntryType = "credit"
)

type Ledger struct {
	ID           string          `json:"id"`
	AccountID    string          `json:"account_id"`
	Amount       int64           `json:"amount"`
	EntryType    LedgerEntryType `json:"entry_type"`
	Reference    string          `json:"reference"`
	BalanceAfter int64           `json:"balance_after"`
	Description  string          `json:"description,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
