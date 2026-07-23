package gateway

import (
	"context"
	transferpb "github.com/Eucastan/eucastanpay/common/proto/transfer"
	transferReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/transfer"
	transferResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/transfer"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/mapper"
)

type TransferGateway struct {
	client transferpb.TransferServiceClient
}

func NewTransferGateway(client transferpb.TransferServiceClient) *TransferGateway {
	return &TransferGateway{
		client: client,
	}
}

func (s *TransferGateway) Transfer(ctx context.Context, input *transferReq.TransferRequest) (*transferResp.TransferResp, error) {
	grpcResp, err := s.client.Transfer(
		ctx,
		mapper.ToProtoTransfer(*input),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToTransferResponse(grpcResp)
	return &resp, nil
}

func (s *TransferGateway) ReverseTransfer(ctx context.Context, userID, originalRef, idemKey string) (*transferResp.MessageResponse, error) {
	grpcResp, err := s.client.ReverseTransfer(
		ctx,
		mapper.ToProtoReverseTransfer(userID, originalRef, idemKey),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToReverseResponse(grpcResp)
	return &resp, nil
}

func (s *TransferGateway) ReconcileAccount(ctx context.Context, accID string, accNo int64) (*transferResp.MessageResponse, error) {
	grpcResp, err := s.client.ReconcileAccount(
		ctx,
		mapper.ToProtoReconciliation(accID, accNo),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToReconcileResponse(grpcResp)
	return &resp, nil
}

func (s *TransferGateway) GetAllTransfers(ctx context.Context) (*transferResp.UserTransferResponse, error) {
	grpcResp, err := s.client.GetAllTransfers(
		ctx,
		mapper.ToProtoListTransfer(),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToListTransferResponse(grpcResp)
	return resp, nil
}

func (s *TransferGateway) GetTransfer(ctx context.Context, transferID string) (*transferResp.TransferResponse, error) {
	grpcResp, err := s.client.GetTransfer(
		ctx,
		mapper.ToProtoTransferByID(transferID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToGetTransferResponse(grpcResp)
	return &resp, nil
}

func (s *TransferGateway) GetTransferByUserID(ctx context.Context, userID string) (*transferResp.TransferResponse, error) {
	grpcResp, err := s.client.GetTransferByUserID(
		ctx,
		mapper.ToProtoTransferByUserID(userID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToGetTransferResponse(grpcResp)
	return &resp, nil
}
