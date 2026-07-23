package service

import (
	"context"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/auth"
	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/common/pkg/security"
	"github.com/Eucastan/eucastanpay/services/admin/config"
	"github.com/Eucastan/eucastanpay/services/admin/internal/domain"
	"github.com/Eucastan/eucastanpay/services/admin/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/admin/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/admin/internal/repository"
	"github.com/sirupsen/logrus"
)

type AdminUseCase struct {
	repo   repository.AdminRepository
	cfg    *config.Config
	logger *logrus.Logger
}

func NewAdminUseCase(repo repository.AdminRepository, cfg *config.Config, logger *logrus.Logger) *AdminUseCase {
	return &AdminUseCase{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (u *AdminUseCase) CreateAdmin(ctx context.Context, input *request.CreateAdminRequest) (*response.AdminResponse, error) {
	admins, err := u.repo.List(ctx, 1, 0)
	if err != nil {
		u.logger.WithError(err).Error("Failed to check existing admins")
		return nil, err
	}

	if len(admins) > 0 {
		return nil, errmessage.ErrBootstrapClosed
	}

	if _, err := u.repo.FindByEmail(ctx, input.Email); err == nil {
		return nil, errmessage.ErrDuplicateEmail
	}

	// hash password
	hash, err := security.GeneratePassHash(input.Password)
	if err != nil {
		return nil, err
	}

	admin := domain.ToAdminDB(hash, input.Role, input)

	if err := u.repo.Create(ctx, admin); err != nil {
		return nil, err
	}

	return response.ToAdminResponse(admin), nil

}

func (u *AdminUseCase) Login(ctx context.Context, input *request.AdminLoginRequest) (*response.AdminLoginResponse, error) {

	admin, err := u.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, errmessage.ErrInvalidCredentials
	}

	u.logger.Info(
		admin.Email,
		admin.TwoFAEnabled,
		admin.TwoFASecret,
		input.TotpCode,
	)

	u.logger.Info("status:", admin.Status)
	u.logger.Info("password hash:", admin.PasswordHash)
	u.logger.Info("input password:", input.Password)

	if admin.Status != domain.StatusActive {
		u.logger.Info("failed at status")
		return nil, errmessage.ErrAccountDisabled
	}

	err = security.IsMatch(admin.PasswordHash, input.Password)
	if err != nil {
		return nil, errmessage.ErrInvalidCredentials
	}

	// Update last login
	now := time.Now()
	admin.LastLoginAt = &now
	if err := u.repo.Update(ctx, admin); err != nil {
		return nil, err
	}

	token, err := auth.GenerateAdminAccessToken(admin.ID, string(admin.Role), u.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.RefreshToken(admin.ID, admin.Email, string(admin.Role), u.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	data := response.ToAdminResponse(admin)

	return &response.AdminLoginResponse{
		Message:      "Login successful",
		Data:         *data,
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (u *AdminUseCase) GetAdminByID(ctx context.Context, id string) (*response.AdminResponse, error) {
	admin, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return response.ToAdminResponse(admin), nil
}

func (u *AdminUseCase) ListAdmins(ctx context.Context, limit, offset int) ([]*response.AdminResponse, error) {
	admins, err := u.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	resp := make([]*response.AdminResponse, 0, len(admins))
	for _, a := range admins {
		resp = append(resp, response.ToAdminResponse(&a))
	}

	return resp, nil
}

func (u *AdminUseCase) LogoutByAdminID(ctx context.Context, adminID string) error {
	return nil
}

func (u *AdminUseCase) UpdateAdmin(ctx context.Context, id string, req *request.UpdateAdminRequest) (*response.AdminResponse, error) {

	admin, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Role != nil {
		admin.Role = domain.AdminRole(*req.Role)
	}

	if req.Status != nil {
		admin.Status = domain.AdminStatus(*req.Status)
	}

	if err := u.repo.Update(ctx, admin); err != nil {
		return nil, err
	}

	return response.ToAdminResponse(admin), nil
}

func (u *AdminUseCase) DeleteAdmin(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
