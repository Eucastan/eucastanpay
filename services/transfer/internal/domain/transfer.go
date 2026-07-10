package domain

import (
	"time"
)

type TransferDirection string

const (
	TransferDir TransferDirection = "transfer"
	ReverseDir  TransferDirection = "reverse"
)

type TransferMode string

const (
	IntraBank  TransferMode = "intraBank"
	InterBank  TransferMode = "interBank"
	OwnAccount TransferMode = "own"
)

type TransferStatus string

const (
	TransferStatusPending   TransferStatus = "pending"
	TransferStatusSuccess   TransferStatus = "success"
	TransferStatusReversing TransferStatus = "reversing"
	TransferStatusReversed  TransferStatus = "reversed"
	TransferStatusFailed    TransferStatus = "failed"
)

type TransferStep string

const (
	StepInitiated TransferStep = "initiated"
	StepDebited   TransferStep = "debited"
	StepCredited  TransferStep = "credited"
	StepCompleted TransferStep = "completed"
)

type Transfer struct {
	ID               string            `db:"id" json:"id"`
	UserID           string            `db:"user_id" json:"user_id"`
	Reference        string            `db:"reference" json:"reference"`
	Step             TransferStep      `db:"step" json:"step"`
	FromAccID        string            `db:"from_account_id" json:"from_account_id"`
	FromAccNo        int64             `db:"from_account_no" json:"from_account_no"`
	ToAccID          string            `db:"to_account_id" json:"to_account_id"`
	ToAccNo          int64             `db:"to_account_no" json:"to_account_no"`
	Amount           int64             `db:"amount" json:"amount"`
	Description      string            `db:"description" json:"description"`
	IdempotencyKey   string            `db:"idempotency_key" json:"idempotency_key"`
	Direction        TransferDirection `db:"direction" json:"direction"` // TRANSFER, REVERSE
	Status           TransferStatus    `db:"status" json:"status"`
	Mode             TransferMode      `db:"mode" json:"mode"`
	ReversalRef      string            `db:"reversal_ref" json:"reversal_ref"`
	IsReversed       bool              `db:"is_reversed" json:"is_reversed"`
	FromBalanceAfter int64             `db:"from_balance_after" json:"from_balance_after"`
	ToBalanceAfter   int64             `db:"to_balance_after" json:"to_balance_after"`
	CreatedAt        time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time         `db:"updated_at" json:"updated_at"`
}
