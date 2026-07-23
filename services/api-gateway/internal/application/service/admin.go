package service

import (
	"context"

	adminReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/admin"
	adminResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/admin"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/gateway"
)

type AdminApplication struct {
	gateway *gateway.AdminGateway
}

func NewAdminApplication(gateway *gateway.AdminGateway) *AdminApplication {
	return &AdminApplication{
		gateway: gateway,
	}
}

func (s *AdminApplication) CreateAdmin(ctx context.Context, req *adminReq.CreateAdminRequest) (*adminResp.AdminResponse, error) {

	return s.gateway.Register(ctx, req)
}

func (s *AdminApplication) Login(ctx context.Context, req *adminReq.AdminLoginRequest) (*adminResp.AdminLoginResponse, error) {

	return s.gateway.Login(ctx, req)
}

func (s *AdminApplication) GetAllAdmins(ctx context.Context, limit, page int) (*adminResp.ListAdminsResponse, error) {

	return s.gateway.GetAllAdmins(ctx, limit, page)
}

func (s *AdminApplication) GetAdmin(ctx context.Context, adminID string) (*adminResp.AdminResponse, error) {

	return s.gateway.GetAdmin(ctx, adminID)
}

func (s *AdminApplication) UpdateAdmin(ctx context.Context, adminID string, req *adminReq.UpdateAdminRequest) (*adminResp.MessageResponse, error) {

	return s.gateway.UpdateAdmin(ctx, adminID, req)

}

func (s *AdminApplication) DeleteAdmin(ctx context.Context, adminID string) (*adminResp.MessageResponse, error) {

	return s.gateway.DeleteAdmin(ctx, adminID)

}
