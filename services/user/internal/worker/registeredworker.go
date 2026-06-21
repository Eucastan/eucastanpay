package worker

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/user/internal/repository"
	"github.com/jackc/pgx/v5"
)

type PublishUserRegistration struct {
	Outbox repository.UserRepository
}

func NewPublishUserRegistration(outbox repository.UserRepository) *PublishUserRegistration {
	return &PublishUserRegistration{
		Outbox: outbox,
	}
}

func (p *PublishUserRegistration) OnUserRegistration(ctx context.Context, u *response.UserResponse) error {
	userData := &events.UserRegisteredEvent{
		BaseEvent: events.NewBaseEvent(ctx, "user-service"),
		UserID:    u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		Timestamp: u.CreatedAt.Unix(),
	}

	return p.Outbox.WithTX(ctx, func(tx pgx.Tx) error {
		return p.Outbox.SaveOutboxEvent(ctx, tx, events.TopicUserRegistered, u.ID, userData)
	})
}
