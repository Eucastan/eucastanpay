package response

import (
	"time"

	"github.com/Eucastan/eucastanpay/services/ledger/internal/domain"
)

type LedgerResponse struct {
	ID           string    `json:"id"`
	AccountID    string    `json:"account_id"`
	Amount       int64     `json:"amount"`
	EntryType    string    `json:"entry_type"`
	Reference    string    `json:"reference"`
	BalanceAfter int64     `json:"balance_after"`
	Description  string    `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToLedgerResponse(ledger *domain.Ledger) *LedgerResponse {
	return &LedgerResponse{
		ID:           ledger.ID,
		AccountID:    ledger.AccountID,
		Amount:       ledger.Amount,
		EntryType:    string(ledger.EntryType),
		Reference:    ledger.Reference,
		BalanceAfter: ledger.BalanceAfter,
		Description:  ledger.Description,
		CreatedAt:    ledger.CreatedAt,
		UpdatedAt:    ledger.UpdatedAt,
	}
}

func ToListLedgerResponse(ledger []domain.Ledger) []LedgerResponse {
	var lists []LedgerResponse

	for _, v := range ledger {
		lists = append(lists, LedgerResponse{
			ID:           v.ID,
			AccountID:    v.AccountID,
			Amount:       v.Amount,
			EntryType:    string(v.EntryType),
			Reference:    v.Reference,
			BalanceAfter: v.BalanceAfter,
			Description:  v.Description,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		})

	}

	return lists
}
