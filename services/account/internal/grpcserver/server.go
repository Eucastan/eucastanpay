package grpcserver

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/proto/account"
	"github.com/Eucastan/eucastanpay/services/account/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/account/internal/repository"
	"github.com/Eucastan/eucastanpay/services/account/internal/usecase"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountServiceServer struct {
	account.UnimplementedAccountServiceServer
	ACC  usecase.AccountUseCase
	Repo repository.AccountRepository
}

func NewAccountServiceServer(acc usecase.AccountUseCase, repo repository.AccountRepository) *AccountServiceServer {
	return &AccountServiceServer{
		ACC:  acc,
		Repo: repo,
	}
}

func (s *AccountServiceServer) Credit(ctx context.Context, req *account.GetCreditRequest) (*account.GetCreditResponse, error) {
	credit := &request.CreditRequest{
		AccountNo: req.AccountNo,
		Amount:    req.Amount,
	}

	err := s.Repo.WithTx(ctx, func(tx pgx.Tx) error {

		if err := s.ACC.Credit(ctx, tx, req.Id, credit); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return &account.GetCreditResponse{
		Message: "Credit processed successfully",
	}, nil
}

func (s *AccountServiceServer) Debit(ctx context.Context, req *account.GetDebitRequest) (*account.GetDebitResponse, error) {
	debit := &request.DebitRequest{
		AccountNo: req.AccountNo,
		Amount:    req.Amount,
	}

	err := s.Repo.WithTx(ctx, func(tx pgx.Tx) error {
		return s.ACC.Debit(ctx, tx, req.Id, debit)
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &account.GetDebitResponse{
		Message: "Debit processed successfully",
	}, nil
}

func (s *AccountServiceServer) GetUserAccount(ctx context.Context, req *account.GetUserAccountRequest) (*account.GetAccountResponse, error) {
	resp, err := s.ACC.GetByAccountIDAndUserID(ctx, req.AccountId, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &account.GetAccountResponse{
		Id:          resp.ID,
		UserId:      resp.UserID,
		AccountNo:   resp.AccountNo,
		Balance:     resp.Balance,
		AccountType: resp.AccountType,
		Currency:    resp.Currency,
	}, nil
}

func (s *AccountServiceServer) GetBalance(ctx context.Context, req *account.GetBalanceRequest) (*account.GetAccountResponse, error) {
	resp, err := s.ACC.GetBalance(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &account.GetAccountResponse{
		Id:          resp.ID,
		UserId:      resp.UserID,
		AccountNo:   resp.AccountNo,
		Balance:     resp.Balance,
		AccountType: resp.AccountType,
		Currency:    resp.Currency,
	}, nil
}

func (s *AccountServiceServer) GetAllAccounts(
	ctx context.Context,
	req *account.ListAccountsRequest,
) (*account.ListAccountsResponse, error) {
	accounts, err := s.ACC.GetAllAccount(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get all accounts")
	}

	resp := make([]*account.Account, 0, len(accounts))
	for _, v := range accounts {
		resp = append(resp, &account.Account{
			AccountId:   v.ID,
			UserId:      v.UserID,
			AccountNo:   v.AccountNo,
			Balance:     v.Balance,
			AccountType: v.AccountType,
			Currency:    v.Currency,
			CreatedAt:   timestamppb.New(v.CreatedAt),
		})
	}

	return &account.ListAccountsResponse{
		Accounts: resp,
	}, nil
}

func (s *AccountServiceServer) ActionOnAccount(ctx context.Context, req *account.ActionRequest) (*account.ActionResponse, error) {
	msg, err := s.ACC.ActionOnAccount(ctx, req.AccountId, req.Status, req.AccountNo)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &account.ActionResponse{
		Message: msg,
	}, nil
}
