package domain

import (
	"encoding/json"
	"time"
)

type AuditRead struct {
	ID            string          `db:"id" json:"id"`
	EventType     string          `db:"event_type" json:"event_type"`
	Service       string          `db:"service" json:"service"`
	CorrelationID string          `db:"correlation_id" json:"correlation_id"`
	CausationID   string          `db:"causation_id" json:"causation_id"`
	Reference     string          `db:"reference" json:"reference"`
	AccountID     string          `db:"account_id" json:"account_id"`
	UserID        string          `db:"user_id" json:"user_id"`
	Amount        int64           `db:"amount" json:"amount"`
	Status        string          `db:"status" json:"status"`
	Payload       json.RawMessage `db:"payload" json:"payload"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
}
