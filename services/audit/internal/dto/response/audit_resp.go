package response

import (
	"encoding/json"
	"time"

	"github.com/Eucastan/eucastanpay/services/audit/internal/domain"
)

type ReadResponse struct {
	Data  []AuditReadResponse `json:"data"`
	Count int                 `json:"count"`
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

func ToAuditReadResponse(read *domain.AuditRead) *AuditReadResponse {
	return &AuditReadResponse{
		ID:            read.ID,
		EventType:     read.EventType,
		Service:       read.Service,
		CorrelationID: read.CorrelationID,
		CausationID:   read.CausationID,
		Reference:     read.Reference,
		AccountID:     read.AccountID,
		UserID:        read.UserID,
		Amount:        read.Amount,
		Status:        read.Status,
		Payload:       read.Payload,
		CreatedAt:     read.CreatedAt,
	}
}
