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

func (s *TransferApplication) Transfer(ctx context.Context, userID, idemKey string, input *transferReq.TransferRequest) (*transferResp.TransferResp, error) {
	return s.gateway.Transfer(ctx, userID, idemKey, input)
}

func (s *TransferApplication) ReverseTransfer(ctx context.Context, originalRef, idemKey string) (*transferResp.TransferResponse, error) {
	return s.gateway.ReverseTransfer(ctx, originalRef, idemKey)
}

func (s *TransferApplication) ReconcileAccount(ctx context.Context, accID string) (*transferResp.ReconciliationResult, error) {
	return s.gateway.ReconcileAccount(ctx, accID)
}

func (s *TransferApplication) GetAllTransfers(ctx context.Context) (*transferResp.UserTransferResponse, error) {
	return s.gateway.GetAllTransfers(ctx)
}

func (s *TransferApplication) GetTransfer(ctx context.Context, transferID string) (*transferResp.TransferResponse, error) {
	return s.gateway.GetTransfer(ctx, transferID)
}

func (s *TransferApplication) GetTransferByRef(ctx context.Context, ref string) (*transferResp.TransferResponse, error) {
	return s.gateway.GetTransferByRef(ctx, ref)
}

func (s *TransferApplication) GetTransferUserID(ctx context.Context, userID string) (*transferResp.TransferResponse, error) {
	return s.gateway.GetTransferByUserID(ctx, userID)
}
