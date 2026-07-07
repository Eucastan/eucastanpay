package events

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EventMetadata struct {
	EventID       string `json:"event_id"`
	CorrelationID string `json:"correlation_id"`
	CausationID   string `json:"causation_id"` // who triggered it
	Version       int    `json:"version"`
	Timestamp     int64  `json:"timestamp"`
}

func NewRootEvent(ctx context.Context) EventMetadata {
	return EventMetadata{
		EventID:       uuid.NewString(),
		CorrelationID: getCorrelationID(ctx),
		CausationID:   "",
		Version:       1,
		Timestamp:     time.Now().Unix(),
	}
}

func NewChildEvent(parent EventMetadata) EventMetadata {
	return EventMetadata{
		EventID:       uuid.NewString(),
		CorrelationID: parent.CorrelationID,
		CausationID:   parent.EventID,
		Version:       parent.Version,
		Timestamp:     time.Now().Unix(),
	}
}

func getCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value("correlation_id").(string); ok && id != "" {
		return id
	}
	return uuid.NewString()
}
