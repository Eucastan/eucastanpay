package repository

import (
	"context"
	"github.com/Eucastan/eucastanpay/services/user/internal/domain"
)

type AuthRepository interface {
	Create(ctx context.Context, auth *domain.Token) error
	FindToken(ctx context.Context, token, tokenType string) (*domain.Token, error)
	FindByUserID(ctx context.Context, userID, tokenType string) (*domain.Token, error)
	UpdateToken(ctx context.Context, id, token, tokenType string) error
	Revoked(ctx context.Context, token string) error
	RevokeAllByUser(ctx context.Context, userID string) error
}
