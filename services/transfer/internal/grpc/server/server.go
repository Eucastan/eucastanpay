package server

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/pkg/grpcstatus"
	transferpb "github.com/Eucastan/eucastanpay/common/proto/transfer"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransferServiceServer struct {
	transferpb.UnimplementedTransferServiceServer
	t usecase.TransferUseCase
}

func NewTransferServiceServer(transfer usecase.TransferUseCase) *TransferServiceServer {
	return &TransferServiceServer{
		t: transfer,
	}
}

func (s *TransferServiceServer) Transfer(ctx context.Context, req *transferpb.TransferRequest) (*transferpb.TransferResponse, error) {
	input := &request.TransferRequest{
		ToAccNo:     req.ToAccountNo,
		Amount:      req.Amount,
		Description: req.Description,
		Mode:        req.Mode,
	}

	resp, err := s.t.TransferFromUser(ctx, req.UserId, req.IdempotencyKey, input)
	if err != nil {
		return nil, grpcstatus.ToTransferStatus(err)
	}

	data := &transferpb.Transfer{
		TransferId:       resp.ID,
		Reference:        resp.Reference,
		Step:             resp.Step,
		FromAccId:        resp.FromAccID,
		FromAccNo:        resp.FromAccNo,
		ToAccId:          resp.ToAccID,
		ToAccNo:          resp.ToAccNo,
		Amount:           resp.Amount,
		Description:      resp.Description,
		IdempotencyKey:   resp.IdempotencyKey,
		Status:           resp.Status,
		Mode:             resp.Mode,
		ReversalRef:      resp.ReversalRef,
		IsReversed:       resp.IsReversed,
		FromBalanceAfter: resp.FromBalanceAfter,
		ToBalanceAfter:   resp.ToBalanceAfter,
		CreatedAt:        timestamppb.New(resp.CreatedAt),
	}

	return &transferpb.TransferResponse{
		Message: "Transfer initiated",
		Resp:    data,
	}, nil
}

func (s *TransferServiceServer) ReverseTransfer(ctx context.Context, req *transferpb.ReverseRequest) (*transferpb.ReverseResponse, error) {
	_, err := s.t.ReverseTransfer(ctx, req.UserId, req.Reference, req.IdempotencyKey)
	if err != nil {
		return nil, grpcstatus.ToTransferStatus(err)
	}

	return &transferpb.ReverseResponse{
		Status: "success",
	}, nil
}

func (s *TransferServiceServer) ReconcileAccount(ctx context.Context, req *transferpb.ReconcileAccountRequest) (*transferpb.ReconcileAccountResponse, error) {
	input := request.ReconciliationRequest{
		AccountNo: req.AccountNo,
	}
	err := s.t.ReconcileAccount(ctx, req.AccountId, &input)
	if err != nil {
		return nil, grpcstatus.ToTransferStatus(err)
	}

	return &transferpb.ReconcileAccountResponse{
		Status: "success",
		Valid:  true,
	}, nil
}

func (s *TransferServiceServer) GetAllTransfers(
	ctx context.Context,
	req *transferpb.ListTransfersRequest,
) (*transferpb.ListTransfersResponse, error) {
	transfers, err := s.t.GetAllTransfers(ctx)
	if err != nil {
		return nil, grpcstatus.ToTransferStatus(err)
	}

	resp := make([]*transferpb.Transfer, 0, len(transfers))
	for _, v := range transfers {
		resp = append(resp, &transferpb.Transfer{
			TransferId:       v.ID,
			Reference:        v.Reference,
			Step:             v.Step,
			FromAccId:        v.FromAccID,
			FromAccNo:        v.FromAccNo,
			ToAccNo:          v.ToAccNo,
			Amount:           v.Amount,
			Description:      v.Description,
			IdempotencyKey:   v.IdempotencyKey,
			Type:             v.Direction,
			Status:           v.Status,
			Mode:             v.Mode,
			ReversalRef:      v.ReversalRef,
			IsReversed:       v.IsReversed,
			FromBalanceAfter: v.FromBalanceAfter,
			ToBalanceAfter:   v.ToBalanceAfter,
			CreatedAt:        timestamppb.New(v.CreatedAt),
		})
	}

	return &transferpb.ListTransfersResponse{
		Transfers: resp,
	}, nil
}

func (s *TransferServiceServer) GetTransfer(
	ctx context.Context,
	req *transferpb.TransferIdRequest,
) (*transferpb.GetTransferResponse, error) {

	resp, err := s.t.GetByID(ctx, req.TransferId)
	if err != nil {
		return nil, grpcstatus.ToTransferStatus(err)
	}

	return &transferpb.GetTransferResponse{
		TransferId:       resp.ID,
		Reference:        resp.Reference,
		Step:             resp.Step,
		FromAccId:        resp.FromAccID,
		FromAccNo:        resp.FromAccNo,
		ToAccId:          resp.ToAccID,
		ToAccNo:          resp.ToAccNo,
		Amount:           resp.Amount,
		Description:      resp.Description,
		IdempotencyKey:   resp.IdempotencyKey,
		Status:           resp.Status,
		Mode:             resp.Mode,
		ReversalRef:      resp.ReversalRef,
		IsReversed:       resp.IsReversed,
		FromBalanceAfter: resp.FromBalanceAfter,
		ToBalanceAfter:   resp.ToBalanceAfter,
		CreatedAt:        timestamppb.New(resp.CreatedAt),
	}, nil
}

func (s *TransferServiceServer) GetTransferByUserID(
	ctx context.Context,
	req *transferpb.UserIdRequest,
) (*transferpb.GetTransferResponse, error) {

	resp, err := s.t.GetByUserID(ctx, req.UserId)
	if err != nil {
		return nil, grpcstatus.ToTransferStatus(err)
	}

	return &transferpb.GetTransferResponse{
		TransferId:       resp.ID,
		Reference:        resp.Reference,
		Step:             resp.Step,
		FromAccId:        resp.FromAccID,
		FromAccNo:        resp.FromAccNo,
		ToAccId:          resp.ToAccID,
		ToAccNo:          resp.ToAccNo,
		Amount:           resp.Amount,
		Description:      resp.Description,
		IdempotencyKey:   resp.IdempotencyKey,
		Status:           resp.Status,
		Mode:             resp.Mode,
		ReversalRef:      resp.ReversalRef,
		IsReversed:       resp.IsReversed,
		FromBalanceAfter: resp.FromBalanceAfter,
		ToBalanceAfter:   resp.ToBalanceAfter,
		CreatedAt:        timestamppb.New(resp.CreatedAt),
	}, nil
}
