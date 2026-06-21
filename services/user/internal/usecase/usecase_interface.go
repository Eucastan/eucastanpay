package usecase

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/user/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
)

type UserUseCaseInterface interface {
	Register(ctx context.Context, input *request.RegisterRequest) (*response.UserResponse, error)
	VerifyEmail(ctx context.Context, token string) error
	Login(ctx context.Context, input *request.LoginRequest) (*response.AuthResponse, error)
	GetAllUsers(ctx context.Context) ([]response.UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*response.UserResponse, error)
	UserCurrentStatus(ctx context.Context, id, status string) (string, error)
	RefreshToken(ctx context.Context, oldToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
	ForgotPassword(ctx context.Context, input *request.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error
}
