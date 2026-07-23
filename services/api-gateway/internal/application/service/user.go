package service

import (
	"context"

	userReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/user"
	userResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/user"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/gateway"
)

type UserApplication struct {
	gateway *gateway.UserGateway
}

func NewUserApplication(gateway *gateway.UserGateway) *UserApplication {
	return &UserApplication{
		gateway: gateway,
	}
}

func (s *UserApplication) Register(ctx context.Context, req *userReq.RegisterRequest) (*userResp.RegisterResponse, error) {

	return s.gateway.Register(ctx, req)
}

func (s *UserApplication) Login(ctx context.Context, req *userReq.LoginRequest) (*userResp.AuthResponse, error) {

	return s.gateway.Login(ctx, req)
}

func (s *UserApplication) GetAllUsers(ctx context.Context) (*userResp.ListUsersResponse, error) {

	return s.gateway.GetAllUsers(ctx)
}

func (s *UserApplication) GetUserByID(ctx context.Context, userID string) (*userResp.UserResponse, error) {

	return s.gateway.GetUserByID(ctx, userID)
}

func (s *UserApplication) UpdateUser(ctx context.Context, userID string, req *userReq.UpdateRequest) (*userResp.MessageResponse, error) {

	return s.gateway.UpdateUser(ctx, userID, req)

}

func (s *UserApplication) DeleteUser(ctx context.Context, userID string) (*userResp.MessageResponse, error) {

	return s.gateway.DeleteUser(ctx, userID)

}

func (s *UserApplication) CreateKYC(ctx context.Context, idNumber, idType string) (*userResp.KYCResponse, error) {

	return s.gateway.CreateKYC(ctx, idNumber, idType)

}

func (s *UserApplication) GetKYC(ctx context.Context, userID string) (*userResp.KYCResponse, error) {

	return s.gateway.GetKYC(ctx, userID)

}

func (s *UserApplication) ApproveKYC(ctx context.Context, userID string) (*userResp.KYCResponse, error) {

	return s.gateway.ApproveKYC(ctx, userID)

}
