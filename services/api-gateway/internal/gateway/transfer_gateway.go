package gateway

import (
	"context"

	transferpb "github.com/Eucastan/eucastanpay/common/proto/transfer"
	transferReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/transfer"
	transferResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/transfer"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/mapper"
	"github.com/sirupsen/logrus"
)

type TransferGateway struct {
	client transferpb.TransferServiceClient
	logger *logrus.Logger
}

func NewTransferGateway(client transferpb.TransferServiceClient, logger *logrus.Logger) *TransferGateway {
	return &TransferGateway{
		client: client,
		logger: logger,
	}
}

func (g *TransferGateway) Transfer(ctx context.Context, userID, idemKey string, input *transferReq.TransferRequest) (*transferResp.TransferResp, error) {
	grpcResp, err := g.client.Transfer(
		ctx,
		mapper.ToProtoTransfer(userID, idemKey, *input),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToTransferResponse(grpcResp)
	return &resp, nil
}

func (g *TransferGateway) ReverseTransfer(ctx context.Context, originalRef, idemKey string) (*transferResp.TransferResponse, error) {
	grpcResp, err := g.client.ReverseTransfer(
		ctx,
		mapper.ToProtoReverseTransfer(originalRef, idemKey),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToReverseResponse(grpcResp)
	return &resp, nil
}

func (g *TransferGateway) ReconcileAccount(ctx context.Context, accID string) (*transferResp.ReconciliationResult, error) {
	grpcResp, err := g.client.ReconcileAccount(
		ctx,
		mapper.ToProtoReconciliation(accID),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToReconcileResponse(grpcResp)
	return &resp, nil
}

func (g *TransferGateway) GetAllTransfers(ctx context.Context) (*transferResp.UserTransferResponse, error) {
	grpcResp, err := g.client.GetAllTransfers(
		ctx,
		mapper.ToProtoListTransfer(),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToListTransferResponse(grpcResp)
	return resp, nil
}

func (g *TransferGateway) GetTransfer(ctx context.Context, transferID string) (*transferResp.TransferResponse, error) {
	grpcResp, err := g.client.GetTransfer(
		ctx,
		mapper.ToProtoTransferByID(transferID),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToGetTransferResponse(grpcResp)
	return &resp, nil
}

func (g *TransferGateway) GetTransferByRef(ctx context.Context, ref string) (*transferResp.TransferResponse, error) {
	grpcResp, err := g.client.GetTransferByRef(
		ctx,
		mapper.ToProtoTransferByRef(ref),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToGetTransferResponse(grpcResp)
	return &resp, nil
}

func (g *TransferGateway) GetTransferByUserID(ctx context.Context, userID string) (*transferResp.TransferResponse, error) {
	grpcResp, err := g.client.GetTransferByUserID(
		ctx,
		mapper.ToProtoTransferByUserID(userID),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToGetTransferResponse(grpcResp)
	return &resp, nil
}
