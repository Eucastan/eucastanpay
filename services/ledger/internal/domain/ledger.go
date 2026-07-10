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
	ID           string          `db:"id" json:"id"`
	UserID       string          `db:"user_id" json:"user_id"`
	AccountID    string          `db:"account_id" json:"account_id"`
	Amount       int64           `db:"amount" json:"amount"`
	EntryType    LedgerEntryType `db:"entry_type" json:"entry_type"`
	Reference    string          `db:"reference" json:"reference"`
	BalanceAfter int64           `db:"balance_after" json:"balance_after"`
	Description  string          `db:"description" json:"description,omitempty"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updated_at"`
}
