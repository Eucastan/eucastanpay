package audit

import (
	"encoding/json"
	"time"
)

type ReadResponse struct {
	Data []AuditReadResponse `json:"data"`
}

type SearchResult struct {
	Count int                 `json:"count"`
	Data  []AuditReadResponse `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type AuditReadResponse struct {
	ID            string          `json:"id"`
	EventType     string          `json:"event_type"`
	Service       string          `json:"service"`
	CorrelationID string          `json:"correlation_id"`
	CausationID   string          `json:"causation_id"`
	Reference     string          `json:"reference"`
	AccountID     string          `json:"account_id"`
	UserID        string          `json:"user_id"`
	Amount        int64           `json:"amount"`
	Status        string          `json:"status"`
	Payload       json.RawMessage `json:"payload"`
	CreatedAt     time.Time       `json:"created_at"`
}
