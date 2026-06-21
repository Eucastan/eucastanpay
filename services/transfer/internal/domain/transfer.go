package domain

import (
	"time"
)

type TransferType string

const (
	TransferTypeDebit   TransferType = "debit"
	TransferTypeCredit  TransferType = "credit"
	TransferTypeReverse TransferType = "reverse"
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
	TransferStatusReverse   TransferStatus = "reverse"
	TransferStatusFailed    TransferStatus = "failed"
)

const (
	StepInitiated = "initiated"
	StepDebited   = "debited"
	StepCredited  = "credited"
)

type Transfer struct {
	ID               string         `json:"id"`
	UserID           string         `json:"user_id"`
	Reference        string         `json:"reference"`
	Step             string         `json:"step"`
	FromAccID        string         `json:"from_acc_id"`
	FromAccNo        int64          `json:"from_acc_no"`
	ToAccID          string         `json:"to_acc_id,omitempty"` // Inter-bank transfer
	ToAccNo          int64          `json:"to_acc_no"`
	Amount           int64          `json:"amount"`
	Description      string         `json:"description"`
	IdempotencyKey   string         `json:"idempotency_key"`
	Type             TransferType   `json:"type"` // DEBIT, CREDIT, REVERSE
	Status           TransferStatus `json:"status"`
	Mode             TransferMode   `json:"mode"`
	ReversalRef      string         `json:"reversal_ref"`
	IsReversed       bool           `json:"is_reversed"`
	FromBalanceAfter int64          `json:"from_balance_after"`
	ToBalanceAfter   int64          `json:"to_balance_after"`
	RecoveryCount    int            `json:"recovery_count"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}
