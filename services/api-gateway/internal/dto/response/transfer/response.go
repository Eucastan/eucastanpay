package transfer

import "time"

type UserTransferResponse struct {
	Message string              `json:"message"`
	Data    []*TransferResponse `json:"data"`
}

type TransferResp struct {
	Message string           `json:"message"`
	Data    TransferResponse `json:"data"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TransferResponse struct {
	ID               string    `json:"id"`
	Reference        string    `json:"reference"`
	Step             string    `json:"step"`
	FromAccID        string    `json:"from_account_id"`
	FromAccNo        int64     `json:"from_account_no"`
	ToAccID          string    `json:"to_account_id"`
	ToAccNo          int64     `json:"to_account_no"`
	Amount           int64     `json:"amount"`
	Description      string    `json:"description"`
	IdempotencyKey   string    `json:"idempotency_key"`
	Direction        string    `json:"type"` // TRANSFER, REVERSE
	Status           string    `json:"status"`
	Mode             string    `json:"mode"`
	ReversalRef      string    `json:"reversal_ref"`
	IsReversed       bool      `json:"is_reversed"`
	FromBalanceAfter int64     `json:"from_balance_after"`
	ToBalanceAfter   int64     `json:"to_balance_after"`
	CreatedAt        time.Time `json:"created_at"`
}
