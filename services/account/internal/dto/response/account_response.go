package response

import (
	"time"

	"github.com/Eucastan/eucastanpay/services/account/internal/domain"
)

type AccountResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	AccountNo   int64     `json:"account_no"`
	Balance     int64     `json:"amount"`
	AccountType string    `json:"account_type"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToAccountResponse(acc *domain.Account) *AccountResponse {
	if acc == nil {
		return nil
	}

	return &AccountResponse{
		ID:          acc.ID,
		UserID:      acc.UserID,
		AccountNo:   acc.AccountNo,
		Balance:     acc.Balance,
		AccountType: string(acc.AccountType),
		Currency:    acc.Currency,
		CreatedAt:   acc.CreatedAt,
	}
}
