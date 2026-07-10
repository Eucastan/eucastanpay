package response

import (
	"time"

	"github.com/Eucastan/eucastanpay/services/notification/internal/domain"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type NotificationResponse struct {
	ID           string                  `json:"id"`
	UserID       string                  `json:"user_id"`
	Title        string                  `json:"title"`
	Message      string                  `json:"message"`
	Channel      domain.Channel          `json:"channel"`
	Type         domain.NotificationType `json:"type"`
	Priority     domain.Priority         `json:"priority"`
	Reference    string                  `json:"reference"`
	Metadata     map[string]any          `json:"metadata"`
	Status       string                  `json:"status"`
	ScheduledFor *time.Time              `json:"scheduled_for"`
	SentAt       *time.Time              `json:"sent_at"`
	CreatedAt    time.Time               `json:"created_at"`
}

type NotificationTemplateResponse struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Subject string         `json:"subject"`
	Body    string         `json:"body"`
	Channel domain.Channel `json:"channel"`
}

func ToNotificationResponse(n []domain.Notification) []NotificationResponse {
	resp := make([]NotificationResponse, 0, len(n))
	for _, v := range n {
		resp = append(resp, NotificationResponse{
			ID:           v.ID,
			UserID:       v.UserID,
			Title:        v.Title,
			Message:      v.Message,
			Channel:      v.Channel,
			Type:         v.Type,
			Priority:     v.Priority,
			Reference:    v.Reference,
			Metadata:     v.Metadata,
			Status:       v.Status,
			ScheduledFor: v.ScheduledFor,
			SentAt:       v.SentAt,
			CreatedAt:    v.CreatedAt,
		})
	}
	return resp
}

func ToNotificationTemplateResponse(n *domain.NotificationTemplate) *NotificationTemplateResponse {
	resp := NotificationTemplateResponse{
		ID:      n.ID,
		Name:    n.Name,
		Subject: n.Subject,
		Body:    n.Body,
		Channel: n.Channel,
	}

	return &resp
}
