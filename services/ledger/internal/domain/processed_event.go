package domain

import (
	"time"
)

type ProcessedEvent struct {
	ID          string    `db:"id" json:"id"`
	EventID     string    `db:"event_id" json:"event_id"`
	Topic       string    `db:"topic" json:"topic"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at"`
}
