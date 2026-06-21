package domain

import (
	"time"
)

type ProcessedEvent struct {
	ID          string    `json:"id"`
	EventID     string    `json:"event_id"`
	Topic       string    `json:"topic"`
	ProcessedAt time.Time `json:"processed_at"`
}
