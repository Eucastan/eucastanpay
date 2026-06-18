package events

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type BaseEvent struct {
	EventID       string `json:"event_id"`
	CorrelationID string `json:"correlation_id"`
	CausationID   string `json:"causation_id"` // who triggered it
	Version       int    `json:"version"`
	Timestamp     int64  `json:"timestamp"`
}

func NewBaseEvent(ctx context.Context, causationID string) BaseEvent {
	return BaseEvent{
		EventID:       uuid.NewString(),
		CorrelationID: getCorrelationID(ctx),
		CausationID:   causationID,
		Version:       1,
		Timestamp:     time.Now().Unix(),
	}
}

func getCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value("correlation_id").(string); ok && id != "" {
		return id
	}
	return uuid.NewString()
}
