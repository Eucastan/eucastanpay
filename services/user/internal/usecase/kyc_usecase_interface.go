package usecase

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/user/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
)

type KYCUseCase interface {
	CreateKYC(ctx context.Context, userID string, input *request.KYCRequest) error
	GetKYC(ctx context.Context, id string) (*response.KYCResponse, error)
	ApproveKYC(ctx context.Context, userID string) error
}
