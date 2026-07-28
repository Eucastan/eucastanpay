package mapper

import (
	transferpb "github.com/Eucastan/eucastanpay/common/proto/transfer"
	transferReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/transfer"
	transferResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/transfer"
)

func ToProtoTransfer(userID, idemKey string, req transferReq.TransferRequest) *transferpb.TransferRequest {
	return &transferpb.TransferRequest{
		UserId:         userID,
		IdempotencyKey: idemKey,
		ToAccountNo:    req.ToAccNo,
		Amount:         req.Amount,
		Description:    req.Description,
		Mode:           req.Mode,
	}
}

func ToProtoReverseTransfer(originalRef, idemKey string) *transferpb.ReverseRequest {
	return &transferpb.ReverseRequest{
		Reference:      originalRef,
		IdempotencyKey: idemKey,
	}
}

func ToProtoReconciliation(accID string) *transferpb.ReconcileAccountRequest {
	return &transferpb.ReconcileAccountRequest{
		AccountId: accID,
	}
}

func ToProtoTransferByRef(ref string) *transferpb.ReferenceRequest {
	return &transferpb.ReferenceRequest{
		Reference: ref,
	}
}

func ToProtoTransferByUserID(userID string) *transferpb.UserIdRequest {
	return &transferpb.UserIdRequest{
		UserId: userID,
	}
}

func ToProtoTransferByID(transferID string) *transferpb.TransferIdRequest {
	return &transferpb.TransferIdRequest{
		TransferId: transferID,
	}
}

func ToProtoListTransfer() *transferpb.ListTransfersRequest {
	return &transferpb.ListTransfersRequest{}
}

func ToReverseResponse(resp *transferpb.ReverseResponse) transferResp.TransferResponse {
	return transferResp.TransferResponse{
		ID:               resp.Data.TransferId,
		Reference:        resp.Data.Reference,
		Step:             resp.Data.Step,
		FromAccID:        resp.Data.FromAccId,
		FromAccNo:        resp.Data.FromAccNo,
		ToAccID:          resp.Data.ToAccId,
		ToAccNo:          resp.Data.ToAccNo,
		Amount:           resp.Data.Amount,
		Description:      resp.Data.Description,
		IdempotencyKey:   resp.Data.IdempotencyKey,
		Direction:        resp.Data.Description,
		Status:           resp.Data.Status,
		Mode:             resp.Data.Mode,
		ReversalRef:      resp.Data.ReversalRef,
		IsReversed:       resp.Data.IsReversed,
		FromBalanceAfter: resp.Data.FromBalanceAfter,
		ToBalanceAfter:   resp.Data.ToBalanceAfter,
		CreatedAt:        resp.Data.CreatedAt.AsTime(),
	}
}

func ToReconcileResponse(resp *transferpb.ReconcileAccountResponse) transferResp.ReconciliationResult {
	return transferResp.ReconciliationResult{
		AccountID:      resp.AccountId,
		AccountBalance: resp.AccountBalance,
		LedgerBalance:  resp.LedgerBalance,
		Difference:     resp.Difference,
		Status:         resp.Status,
		Reason:         resp.Reason,
		ReconciledAt:   resp.ReconciledAt.AsTime(),
	}
}

func ToGetTransferResponse(resp *transferpb.GetTransferResponse) transferResp.TransferResponse {

	return transferResp.TransferResponse{
		ID:               resp.TransferId,
		Reference:        resp.Reference,
		Step:             resp.Step,
		FromAccID:        resp.FromAccId,
		FromAccNo:        resp.FromAccNo,
		ToAccID:          resp.ToAccId,
		ToAccNo:          resp.ToAccNo,
		Amount:           resp.Amount,
		Description:      resp.Description,
		IdempotencyKey:   resp.IdempotencyKey,
		Direction:        resp.Description,
		Status:           resp.Status,
		Mode:             resp.Mode,
		ReversalRef:      resp.ReversalRef,
		IsReversed:       resp.IsReversed,
		FromBalanceAfter: resp.FromBalanceAfter,
		ToBalanceAfter:   resp.ToBalanceAfter,
		CreatedAt:        resp.CreatedAt.AsTime(),
	}

}

func ToTransferResponse(resp *transferpb.TransferResponse) transferResp.TransferResp {

	data := transferResp.TransferResponse{
		ID:               resp.Resp.TransferId,
		Reference:        resp.Resp.Reference,
		Step:             resp.Resp.Step,
		FromAccID:        resp.Resp.FromAccId,
		FromAccNo:        resp.Resp.FromAccNo,
		ToAccID:          resp.Resp.ToAccId,
		ToAccNo:          resp.Resp.ToAccNo,
		Amount:           resp.Resp.Amount,
		Description:      resp.Resp.Description,
		IdempotencyKey:   resp.Resp.IdempotencyKey,
		Direction:        resp.Resp.Description,
		Status:           resp.Resp.Status,
		Mode:             resp.Resp.Mode,
		ReversalRef:      resp.Resp.ReversalRef,
		IsReversed:       resp.Resp.IsReversed,
		FromBalanceAfter: resp.Resp.FromBalanceAfter,
		ToBalanceAfter:   resp.Resp.ToBalanceAfter,
		CreatedAt:        resp.Resp.CreatedAt.AsTime(),
	}

	return transferResp.TransferResp{
		Message: resp.Message,
		Data:    data,
	}
}

func ToListTransferResponse(resp *transferpb.ListTransfersResponse) *transferResp.UserTransferResponse {

	data := make([]*transferResp.TransferResponse, 0, len(resp.Transfers))
	for _, r := range resp.Transfers {
		data = append(data, &transferResp.TransferResponse{
			ID:               r.TransferId,
			Reference:        r.Reference,
			Step:             r.Step,
			FromAccID:        r.FromAccId,
			FromAccNo:        r.FromAccNo,
			ToAccID:          r.ToAccId,
			ToAccNo:          r.ToAccNo,
			Amount:           r.Amount,
			Description:      r.Description,
			IdempotencyKey:   r.IdempotencyKey,
			Direction:        r.Description,
			Status:           r.Status,
			Mode:             r.Mode,
			ReversalRef:      r.ReversalRef,
			IsReversed:       r.IsReversed,
			FromBalanceAfter: r.FromBalanceAfter,
			ToBalanceAfter:   r.ToBalanceAfter,
			CreatedAt:        r.CreatedAt.AsTime(),
		})
	}

	return &transferResp.UserTransferResponse{
		Data: data,
	}
}
