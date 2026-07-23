package audit

import "time"

type Filter struct {
	CorrelationID string
	Reference     string
	EventType     string
	MinAmount     int64
	MaxAmount     int64
	FromDate      *time.Time
	ToDate        *time.Time
	Limit         int
	Offset        int
}

type AuditIdRequest struct {
	AuditID string `json:"id"`
}
