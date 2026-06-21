package server

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/proto/transfer"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransferServiceServer struct {
	transfer.UnimplementedTransferServiceServer
	Transfer usecase.TransferUseCase
}

func NewTransferServiceServer(transfer usecase.TransferUseCase) *TransferServiceServer {
	return &TransferServiceServer{
		Transfer: transfer,
	}
}

func (s *TransferServiceServer) ReverseTransfer(ctx context.Context, req *transfer.ReverseRequest) (*transfer.ReverseResponse, error) {
	_, err := s.Transfer.ReverseTransfer(ctx, req.UserId, req.Reference, req.IdempotencyKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to initiate transaction")
	}

	return &transfer.ReverseResponse{
		Status: "success",
	}, nil
}

func (s *TransferServiceServer) ReconcileAccount(ctx context.Context, req *transfer.ReconcileAccountRequest) (*transfer.ReconcileAccountResponse, error) {
	err := s.Transfer.ReconcileAccount(ctx, req.AccountId, req.AccountNo)
	if err != nil {
		return nil, status.Error(codes.Internal, "account reconciliation failed")
	}

	return &transfer.ReconcileAccountResponse{
		Status: "success",
		Valid:  true,
	}, nil
}

func (s *TransferServiceServer) GetAllTransfers(
	ctx context.Context,
	req *transfer.ListTransfersRequest,
) (*transfer.ListTransfersResponse, error) {
	transfers, err := s.Transfer.GetAllTransfers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get all transfers")
	}

	resp := make([]*transfer.Transfer, 0, len(transfers))
	for _, v := range transfers {
		resp = append(resp, &transfer.Transfer{
			TransferId:       v.ID,
			Reference:        v.Reference,
			Step:             v.Step,
			FromAccId:        v.FromAccID,
			FromAccNo:        v.FromAccNo,
			ToAccId:          v.ToAccID,
			ToAccNo:          v.ToAccNo,
			Amount:           v.Amount,
			Description:      v.Description,
			IdempotencyKey:   v.IdempotencyKey,
			Type:             v.Type,
			Status:           v.Status,
			Mode:             v.Mode,
			ReversalRef:      v.ReversalRef,
			IsReversed:       v.IsReversed,
			FromBalanceAfter: v.FromBalanceAfter,
			ToBalanceAfter:   v.ToBalanceAfter,
			CreatedAt:        timestamppb.New(v.CreatedAt),
		})
	}

	return &transfer.ListTransfersResponse{
		Transfers: resp,
	}, nil
}
