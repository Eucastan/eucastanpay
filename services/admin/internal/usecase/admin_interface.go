package usecase

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/admin/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/admin/internal/dto/response"
)

type AdminUseCase interface {
	CreateAdmin(ctx context.Context, input *request.CreateAdminRequest) (*response.AdminResponse, error)
	Login(ctx context.Context, input *request.AdminLoginRequest) (*response.AdminLoginResponse, error)
	GetAdminByID(ctx context.Context, id string) (*response.AdminResponse, error)
	ListAdmins(ctx context.Context, limit, offset int) ([]*response.AdminResponse, error)
	LogoutByAdminID(ctx context.Context, adminID string) error
	UpdateAdmin(ctx context.Context, id string, input *request.UpdateAdminRequest) (*response.AdminResponse, error)
	DeleteAdmin(ctx context.Context, id string) error
}
