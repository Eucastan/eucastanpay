package service

import (
	"context"

	transferReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/transfer"
	transferResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/transfer"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/gateway"
)

type TransferApplication struct {
	gateway *gateway.TransferGateway
}

func NewTransferApplication(gateway *gateway.TransferGateway) *TransferApplication {
	return &TransferApplication{
		gateway: gateway,
	}
}

func (s *TransferApplication) Transfer(ctx context.Context, input *transferReq.TransferRequest) (*transferResp.TransferResp, error) {
	return s.gateway.Transfer(ctx, input)
}

func (s *TransferApplication) ReverseTransfer(ctx context.Context, userID, originalRef, idemKey string) (*transferResp.MessageResponse, error) {
	return s.gateway.ReverseTransfer(ctx, userID, originalRef, idemKey)
}

func (s *TransferApplication) ReconcileAccount(ctx context.Context, accID string, accNo int64) (*transferResp.MessageResponse, error) {
	return s.gateway.ReconcileAccount(ctx, accID, accNo)
}

func (s *TransferApplication) GetAllTransfers(ctx context.Context) (*transferResp.UserTransferResponse, error) {
	return s.gateway.GetAllTransfers(ctx)
}

func (s *TransferApplication) GetTransfer(ctx context.Context, transferID string) (*transferResp.TransferResponse, error) {
	return s.gateway.GetTransfer(ctx, transferID)
}

func (s *TransferApplication) GetTransferUserID(ctx context.Context, userID string) (*transferResp.TransferResponse, error) {
	return s.gateway.GetTransferByUserID(ctx, userID)
}
