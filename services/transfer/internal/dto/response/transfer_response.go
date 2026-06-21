package response

import (
	"time"

	"github.com/Eucastan/eucastanpay/services/transfer/internal/domain"
)

type TransferResponse struct {
	ID               string    `json:"id"`
	Reference        string    `json:"reference"`
	Step             string    `json:"step"`
	FromAccID        string    `json:"from_account_id"`
	FromAccNo        int64     `json:"from_account_no"`
	ToAccID          string    `json:"to_account_id,omitempty"` // Inter-bank transfer
	ToAccNo          int64     `json:"to_account_no"`
	Amount           int64     `json:"amount"`
	Description      string    `json:"description"`
	IdempotencyKey   string    `json:"idempotency_key"`
	Type             string    `json:"type"` // DEBIT, CREDIT, REVERSE
	Status           string    `json:"status"`
	Mode             string    `json:"mode"`
	ReversalRef      string    `json:"reversal_ref"`
	IsReversed       bool      `json:"is_reversed"`
	FromBalanceAfter int64     `json:"from_balance_after"`
	ToBalanceAfter   int64     `json:"to_balance_after"`
	CreatedAt        time.Time `json:"created_at"`
}

func ToTransferResponse(t *domain.Transfer) TransferResponse {
	return TransferResponse{
		ID:               t.ID,
		Reference:        t.Reference,
		Step:             t.Step,
		FromAccID:        t.FromAccID,
		FromAccNo:        t.FromAccNo,
		ToAccID:          t.ToAccID,
		ToAccNo:          t.ToAccNo,
		Amount:           t.Amount,
		Description:      t.Description,
		IdempotencyKey:   t.IdempotencyKey,
		Type:             string(t.Type),
		Status:           string(t.Status),
		Mode:             string(t.Mode),
		ReversalRef:      t.ReversalRef,
		IsReversed:       t.IsReversed,
		FromBalanceAfter: t.FromBalanceAfter,
		ToBalanceAfter:   t.ToBalanceAfter,
		CreatedAt:        t.CreatedAt,
	}
}
