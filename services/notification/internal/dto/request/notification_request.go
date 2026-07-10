package request

import (
	"github.com/Eucastan/eucastanpay/services/notification/internal/domain"
)

type NotificationRequest struct {
	Title     string                  `json:"title"`
	Message   string                  `json:"message"`
	Channel   domain.Channel          `json:"channel"`
	Type      domain.NotificationType `json:"type"`
	Priority  domain.Priority         `json:"priority"`
	Reference string                  `json:"reference,omitempty"`
	Metadata  map[string]any          `json:"metadata,omitempty"`
}

type NotificationTemplateRequest struct {
	Name    string         `json:"name"`
	Subject string         `json:"subject"`
	Body    string         `json:"body"`
	Channel domain.Channel `json:"channel"`
}
