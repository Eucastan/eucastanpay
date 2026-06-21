package usecase

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/transfer/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/dto/response"
)

type TransferUseCase interface {
	GetAllTransfers(ctx context.Context) ([]response.TransferResponse, error)
	GetByID(ctx context.Context, id string) (*response.TransferResponse, error)
	Transfer(ctx context.Context, userID string, idemKey string, input *request.TransferRequest) (*response.TransferResponse, error)
	TransferFromUser(ctx context.Context, userID string, idemKey string, input *request.TransferRequest) (*response.TransferResponse, error)
	ReverseTransfer(ctx context.Context, userID, originalRef, idemKey string) (*response.TransferResponse, error)
	ReconcileAccount(ctx context.Context, accID string, accNo int64) error
}
