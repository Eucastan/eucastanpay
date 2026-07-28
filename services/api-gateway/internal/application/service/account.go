package service

import (
	"context"

	accountReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/account"
	accountResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/account"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/gateway"
)

type AccountApplication struct {
	gateway *gateway.AccountGateway
}

func NewAccountApplication(gateway *gateway.AccountGateway) *AccountApplication {
	return &AccountApplication{
		gateway: gateway,
	}
}

func (s *AccountApplication) Deposit(ctx context.Context, accID string, input *accountReq.DepositRequest) (*accountResp.MessageResponse, error) {
	return s.gateway.Deposit(ctx, accID, input)
}

func (s *AccountApplication) WithDraw(ctx context.Context, accID string, input *accountReq.DepositRequest) (*accountResp.MessageResponse, error) {
	return s.gateway.WithDraw(ctx, accID, input)
}

func (s *AccountApplication) GetBalance(ctx context.Context, accID, userID string) (*accountResp.AccountResponse, error) {
	return s.gateway.GetBalance(
		ctx, accID, userID,
	)
}

func (s *AccountApplication) GetUserAccount(ctx context.Context, userID string) (*accountResp.AccountResponse, error) {
	return s.gateway.GetUserAccount(
		ctx,
		userID,
	)
}

func (s *AccountApplication) GetAllAccounts(ctx context.Context) ([]*accountResp.AccountResponse, error) {
	return s.gateway.GetAllAccounts(ctx)
}

func (s *AccountApplication) ActionOnAccount(ctx context.Context, accID string, input *accountReq.ActionRequest) (*accountResp.MessageResponse, error) {
	return s.gateway.ActionOnAccount(
		ctx, accID, input,
	)
}

func (s *AccountApplication) DeleteAccount(ctx context.Context, accID string) (*accountResp.MessageResponse, error) {
	return s.gateway.DeleteAccount(
		ctx,
		accID,
	)
}
