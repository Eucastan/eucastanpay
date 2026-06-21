package domain

import (
	"time"
)

type AuditLog struct {
	ID            string                 `json:"id"`
	EventType     string                 `json:"event_type"`
	CorrelationID string                 `json:"correlation_id"`
	CausationID   string                 `json:"causation_id"`
	Reference     string                 `json:"reference"`
	Payload       map[string]interface{} `json:"payload"`
	CreatedAt     time.Time              `json:"created_at"`
}
