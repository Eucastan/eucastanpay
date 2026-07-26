package grpcserver

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	"github.com/Eucastan/eucastanpay/common/pkg/grpcstatus"
	accountpb "github.com/Eucastan/eucastanpay/common/proto/account"
	"github.com/Eucastan/eucastanpay/services/account/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/account/internal/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountServiceServer struct {
	accountpb.UnimplementedAccountServiceServer
	ACC usecase.AccountUseCase
}

func NewAccountServiceServer(acc usecase.AccountUseCase) *AccountServiceServer {
	return &AccountServiceServer{
		ACC: acc,
	}
}

func (s *AccountServiceServer) Deposit(ctx context.Context, req *accountpb.DepositRequest) (*accountpb.ActionResponse, error) {
	_, err := interceptor.RequireAdminRole(ctx, "super_admin", "admin")
	if err != nil {
		return nil, err
	}

	credit := &request.DepositRequest{
		AccountNo: req.AccountNo,
		Amount:    req.Amount,
		Currency:  req.Currency,
	}

	if err := s.ACC.DepositAccount(ctx, req.AccountId, credit); err != nil {
		return nil, grpcstatus.ToAccountStatus(err)
	}

	return &accountpb.ActionResponse{
		Message: "Credit processed successfully",
	}, nil
}

func (s *AccountServiceServer) Withdraw(ctx context.Context, req *accountpb.WithDrawalRequest) (*accountpb.ActionResponse, error) {
	debit := &request.DepositRequest{
		AccountNo: req.AccountNo,
		Amount:    req.Amount,
		Currency:  req.Currency,
	}

	if err := s.ACC.WithDrawal(ctx, req.AccountId, debit); err != nil {

		return nil, grpcstatus.ToAccountStatus(err)
	}

	return &accountpb.ActionResponse{
		Message: "Debit processed successfully",
	}, nil
}

func (s *AccountServiceServer) GetUserAccount(ctx context.Context, req *accountpb.GetUserAccountRequest) (*accountpb.GetAccountResponse, error) {
	resp, err := s.ACC.GetByUserID(ctx, req.UserId)
	if err != nil {
		return nil, grpcstatus.ToAccountStatus(err)
	}

	return &accountpb.GetAccountResponse{
		AccountId:   resp.ID,
		UserId:      resp.UserID,
		Email:       resp.Email,
		AccountNo:   resp.AccountNo,
		Balance:     resp.Balance,
		AccountType: resp.AccountType,
		Currency:    resp.Currency,
	}, nil
}

func (s AccountServiceServer) ResolveAccount(ctx context.Context, req *accountpb.ConfirmAccountRequest) (*accountpb.ConfirmAccountResponse, error) {
	resp, err := s.ACC.ConfirmSenderAndReceiver(ctx, req.FromAccountNo, req.ToAccountNo)
	if err != nil {
		return nil, grpcstatus.ToAccountStatus(err)
	}

	return &accountpb.ConfirmAccountResponse{
		FromAccountId: resp.FromAccID,
		ToAccountId:   resp.ToAccID,
		FromUserId:    resp.FromUserID,
		ToUserId:      resp.ToUserID,
		FromEmail:     resp.FromEmail,
		ToEmail:       resp.ToEmail,
		FromBalance:   resp.FromBalance,
		ToBalance:     resp.ToBalance,
		FromStatus:    resp.FromStatus,
		ToStatus:      resp.ToStatus,
	}, nil
}

func (s *AccountServiceServer) ReconcileBalance(ctx context.Context, req *accountpb.BalanceRequest) (*accountpb.GetAccountResponse, error) {

	resp, err := s.ACC.GetBalance(ctx, req.AccountId, req.UserId)
	if err != nil {
		return nil, grpcstatus.ToAccountStatus(err)
	}

	return &accountpb.GetAccountResponse{
		AccountId:   resp.ID,
		UserId:      resp.UserID,
		Email:       resp.Email,
		AccountNo:   resp.AccountNo,
		Balance:     resp.Balance,
		AccountType: resp.AccountType,
		Currency:    resp.Currency,
	}, nil
}

func (s *AccountServiceServer) GetBalance(ctx context.Context, req *accountpb.GetBalanceRequest) (*accountpb.GetAccountResponse, error) {

	user, err := interceptor.RequireUser(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := s.ACC.GetBalance(ctx, req.AccountId, user.UserID)
	if err != nil {
		return nil, grpcstatus.ToAccountStatus(err)
	}

	return &accountpb.GetAccountResponse{
		AccountId:   resp.ID,
		UserId:      resp.UserID,
		Email:       resp.Email,
		AccountNo:   resp.AccountNo,
		Balance:     resp.Balance,
		AccountType: resp.AccountType,
		Currency:    resp.Currency,
	}, nil
}

func (s *AccountServiceServer) GetAllAccounts(
	ctx context.Context,
	req *accountpb.ListAccountsRequest,
) (*accountpb.ListAccountsResponse, error) {
	_, err := interceptor.RequireAdminRole(ctx, "super_admin", "admin")
	if err != nil {
		return nil, err
	}

	accounts, err := s.ACC.GetAllAccount(ctx)
	if err != nil {
		return nil, grpcstatus.ToAccountStatus(err)
	}

	resp := make([]*accountpb.Account, 0, len(accounts))
	for _, v := range accounts {
		resp = append(resp, &accountpb.Account{
			AccountId:   v.ID,
			UserId:      v.UserID,
			Email:       v.Email,
			AccountNo:   v.AccountNo,
			Balance:     v.Balance,
			AccountType: v.AccountType,
			Currency:    v.Currency,
			CreatedAt:   timestamppb.New(v.CreatedAt),
		})
	}

	return &accountpb.ListAccountsResponse{
		Accounts: resp,
	}, nil
}

func (s *AccountServiceServer) ActionOnAccount(ctx context.Context, req *accountpb.ActionRequest) (*accountpb.ActionResponse, error) {
	_, err := interceptor.RequireAdminRole(ctx, "super_admin", "admin")
	if err != nil {
		return nil, err
	}

	msg, err := s.ACC.ActionOnAccount(ctx, req.AccountId, req.Status, req.AccountNo)
	if err != nil {
		return nil, grpcstatus.ToAccountStatus(err)
	}

	return &accountpb.ActionResponse{
		Message: msg,
	}, nil
}

func (s AccountServiceServer) Delete(ctx context.Context, req *accountpb.DeleteRequest) (*accountpb.ActionResponse, error) {
	_, err := interceptor.RequireAdminRole(ctx, "super_admin", "admin")
	if err != nil {
		return nil, err
	}

	if err := s.ACC.DeleteAccount(ctx, req.AccountId); err != nil {
		return nil, grpcstatus.ToAccountStatus(err)
	}

	return &accountpb.ActionResponse{
		Message: "Deleted successsfully",
	}, nil
}
