package healthcheck

import "time"

type Status string

const (
	Healthy   Status = "healthy"
	Degraded  Status = "degraded"
	Unhealthy Status = "unhealthy"
)

type Component struct {
	Name     string      `json:"name"`
	Status   Status      `json:"status"`
	Duration string      `json:"duration"`
	Details  interface{} `json:"details,omitempty"`
	Error    string      `json:"error,omitempty"`
}

type Summary struct {
	Healthy   int `json:"healthy"`
	Degraded  int `json:"degraded"`
	Unhealthy int `json:"unhealthy"`
}

type Report struct {
	Service    string      `json:"service"`
	Version    string      `json:"version"`
	Status     Status      `json:"status"`
	StartedAt  time.Time   `json:"started_at"`
	Uptime     string      `json:"uptime"`
	Timestamp  time.Time   `json:"timestamp"`
	Summary    Summary     `json:"summary"`
	Components []Component `json:"components"`
}
