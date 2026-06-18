package healthcheck

import (
	"time"
)

type ComponentStatus string

const (
	StatusHealthy   ComponentStatus = "healthy"
	StatusDegraded  ComponentStatus = "degraded"
	StatusUnhealthy ComponentStatus = "unhealthy"
)

type HealthReport struct {
	Status     string                 `json:"status"`
	Service    string                 `json:"service"`
	Version    string                 `json:"version"`
	Uptime     string                 `json:"uptime"`
	Timestamp  time.Time              `json:"timestamp"`
	Components map[string]Component   `json:"components"`
	Summary    map[string]interface{} `json:"summary"`
}

type Component struct {
	Status  ComponentStatus `json:"status"`
	Details interface{}     `json:"details,omitempty"`
	Error   string          `json:"error,omitempty"`
}
