package eventhandler

import (
	"strings"
	"time"
	"encoding/json"

	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/services/audit/internal/domain"
	"github.com/google/uuid"
)

// transform converts raw event payload into structured AuditRead for querying
func transformToRead(topic string, payload map[string]interface{}) *domain.AuditRead {
	now := time.Now()

	read := &domain.AuditRead{
		ID:        uuid.NewString(),
		EventType: topic,
		Service:   getServiceFromTopic(topic),
		CreatedAt: now,
	}

	payloadJSON, _ := json.Marshal(payload)
    read.Payload = payloadJSON

	switch topic {
	case events.TopicUserRegistered:
		read.UserID = getString(payload, "id")
		read.Reference = getString(payload, "id")
		read.Status = "SUCCESS"

	case events.TopicAccountCreated:
		read.AccountID = getString(payload, "id")
		read.UserID = getString(payload, "user_id")
		read.Reference = getString(payload, "id")
		read.Status = "SUCCESS"

	case events.TopicTransferInitiated, events.TopicTransferCompleted:
		read.Reference = getString(payload, "reference")
		read.UserID = getString(payload, "user_id")
		read.AccountID = getString(payload, "from_account_id")
		read.Amount = getInt64(payload, "amount")
		read.Status = "SUCCESS"

	case events.TopicTransferFailed:
		read.Reference = getString(payload, "reference")
		read.UserID = getString(payload, "user_id")
		read.Amount = getInt64(payload, "amount")
		read.Status = "FAILED"

	case events.TopicDebitRequested, events.TopicDebitCompleted:
		read.Reference = getString(payload, "reference")
		read.AccountID = getString(payload, "from_account_id")
		read.Amount = getInt64(payload, "amount")
		read.Status = "SUCCESS"

	case events.TopicCreditRequested, events.TopicCreditCompleted:
		read.Reference = getString(payload, "reference")
		read.AccountID = getString(payload, "to_account_id")
		read.Amount = getInt64(payload, "amount")
		read.Status = "SUCCESS"

	case events.TopicLedgerCreated:
		read.Reference = getString(payload, "reference")
		read.AccountID = getString(payload, "account_id")
		read.Amount = getInt64(payload, "amount")
		read.Status = "SUCCESS"

	case events.TopicDebitFailed, events.TopicCreditFailed:
		read.Reference = getString(payload, "reference")
		read.Amount = getInt64(payload, "amount")
		read.Status = "FAILED"

	default:
		read.Status = "PROCESSED"
	}

	// Try to extract common fields
	if read.Reference == "" {
		read.Reference = getString(payload, "reference")
	}
	if read.CorrelationID == "" {
		read.CorrelationID = getString(payload, "correlation_id")
	}

	return read
}

// Helper to determine service name from topic
func getServiceFromTopic(topic string) string {
	switch {
	case topic == events.TopicUserRegistered || topic == events.TopicUserRegistrationFailed:
		return "user-service"
	case topic == events.TopicAccountCreated || topic == events.TopicCreateAccFailed:
		return "account-service"
	case contains(topic, "transfer"):
		return "transfer-service"
	case contains(topic, "debit") || contains(topic, "credit"):
		return "account-service"
	case contains(topic, "ledger"):
		return "ledger-service"
	default:
		return "unknown-service"
	}
}

func contains(s, substr string) bool {
	// return len(s) >= len(substr) && s[:len(substr)] == substr
	return strings.Contains(s, substr)
}
