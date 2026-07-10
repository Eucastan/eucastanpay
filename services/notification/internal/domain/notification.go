package domain

import (
	"time"
)

type Channel string
type NotificationType string
type Priority string

const (
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelPush  Channel = "push"
	ChannelInApp Channel = "in_app"

	NotificationTypeTransaction NotificationType = "transaction"
	NotificationTypeSecurity    NotificationType = "security"
	NotificationTypeAccount     NotificationType = "account"
	NotificationTypeMarketing   NotificationType = "marketing"
	NotificationTypeSystem      NotificationType = "system"

	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

type Notification struct {
	ID           string           `db:"id" json:"id"`
	UserID       string           `db:"user_id" json:"user_id"`
	Title        string           `db:"title" json:"title"`
	Message      string           `db:"message" json:"message"`
	Channel      Channel          `db:"channel" json:"channel"`
	Type         NotificationType `db:"type" json:"type"`
	Priority     Priority         `db:"priority" json:"priority"`
	Reference    string           `db:"reference" json:"reference,omitempty"` // e.g., transaction reference
	Metadata     map[string]any   `db:"metadata" json:"metadata,omitempty"`   // Extra data (amount, account no, etc.)
	Status       string           `db:"status" json:"status"`                 // pending, sent, failed, delivered
	ScheduledFor *time.Time       `db:"scheduled_for" json:"scheduled_for,omitempty"`
	SentAt       *time.Time       `db:"sent_at" json:"sent_at,omitempty"`
	CreatedAt    time.Time        `db:"created_at" json:"created_at"`
}

type NotificationTemplate struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Subject   string    `db:"subject" json:"subject"`
	Body      string    `db:"body" json:"body"` // Can contain placeholders {{amount}} {{account}}
	Channel   Channel   `db:"channel" json:"channel"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
