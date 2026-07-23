package mapper

import (
	accountpb "github.com/Eucastan/eucastanpay/common/proto/account"

	accountreq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/account"

	accountresp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/account"
)

func ToProtoGetBalance(req accountreq.GetBalanceRequest) *accountpb.GetBalanceRequest {
	return &accountpb.GetBalanceRequest{
		AccountId: req.AccountID,
		UserId:    req.UserID,
	}
}

func ToProtoDeposit(req accountreq.DepositRequest) *accountpb.DepositRequest {
	return &accountpb.DepositRequest{
		AccountId: req.AccountID,
		AccountNo: req.AccountNo,
		Amount:    req.Amount,
		Currency:  req.Currency,
	}
}

func ToProtoWithdraw(req accountreq.DepositRequest) *accountpb.WithDrawalRequest {
	return &accountpb.WithDrawalRequest{
		AccountId: req.AccountID,
		AccountNo: req.AccountNo,
		Amount:    req.Amount,
		Currency:  req.Currency,
	}
}

func ToProtoGetUserAccount(userID string) *accountpb.GetUserAccountRequest {
	return &accountpb.GetUserAccountRequest{
		UserId: userID,
	}
}

func ToProtoListAccount() *accountpb.ListAccountsRequest {
	return &accountpb.ListAccountsRequest{}
}

func ToProtoDeleteRequest(accID string) *accountpb.DeleteRequest {
	return &accountpb.DeleteRequest{
		AccountId: accID,
	}
}

func ToProtoActionRequest(resp *accountreq.ActionRequest) *accountpb.ActionRequest {
	return &accountpb.ActionRequest{
		AccountId: resp.AccountID,
		Status:    resp.Status,
		AccountNo: resp.AccountNo,
	}
}

func ToActionResponse(resp *accountpb.ActionResponse) accountresp.MessageResponse {
	return accountresp.MessageResponse{
		Message: resp.Message,
	}
}

func ToGetAccountResponse(acc *accountpb.GetAccountResponse) accountresp.AccountResponse {

	return accountresp.AccountResponse{
		ID:          acc.AccountId,
		UserID:      acc.UserId,
		Email:       acc.Email,
		AccountNo:   acc.AccountNo,
		Balance:     acc.Balance,
		AccountType: string(acc.AccountType),
		Currency:    acc.Currency,
		Status:      string(acc.Status),
		CreatedAt:   acc.CreatedAt.AsTime(),
		UpdatedAt:   acc.UpdatedAt.AsTime(),
	}
}

func ToListAccountResponse(req *accountpb.ListAccountsResponse) []*accountresp.AccountResponse {
	data := make([]*accountresp.AccountResponse, 0, len(req.Accounts))
	for _, r := range req.Accounts {
		data = append(data, &accountresp.AccountResponse{
			ID:          r.AccountId,
			UserID:      r.UserId,
			Email:       r.Email,
			AccountNo:   r.AccountNo,
			Balance:     r.Balance,
			AccountType: r.AccountType,
			Currency:    r.Currency,
			Status:      r.Status,
			CreatedAt:   r.CreatedAt.AsTime(),
			UpdatedAt:   r.UpdatedAt.AsTime(),
		})
	}
	return data
}
