package gateway

import (
	"context"
	accountpb "github.com/Eucastan/eucastanpay/common/proto/account"
	accountReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/account"
	accountResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/account"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/mapper"
)

type AccountGateway struct {
	client accountpb.AccountServiceClient
}

func NewAccountGateway(client accountpb.AccountServiceClient) *AccountGateway {
	return &AccountGateway{
		client: client,
	}
}

func (g *AccountGateway) Deposit(ctx context.Context, input *accountReq.DepositRequest) (*accountResp.MessageResponse, error) {
	grpcResp, err := g.client.Deposit(
		ctx,
		mapper.ToProtoDeposit(*input),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToActionResponse(grpcResp)
	return &resp, nil
}

func (g *AccountGateway) WithDraw(ctx context.Context, input *accountReq.DepositRequest) (*accountResp.MessageResponse, error) {
	grpcResp, err := g.client.Withdraw(
		ctx,
		mapper.ToProtoWithdraw(*input),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToActionResponse(grpcResp)
	return &resp, nil
}

func (g *AccountGateway) GetBalance(ctx context.Context, input *accountReq.GetBalanceRequest) (*accountResp.AccountResponse, error) {
	grpcResp, err := g.client.GetBalance(
		ctx,
		mapper.ToProtoGetBalance(*input),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToGetAccountResponse(grpcResp)
	return &resp, nil
}

func (g *AccountGateway) GetUserAccount(ctx context.Context, userID string) (*accountResp.AccountResponse, error) {
	grpcResp, err := g.client.GetUserAccount(
		ctx,
		mapper.ToProtoGetUserAccount(userID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToGetAccountResponse(grpcResp)
	return &resp, nil
}

func (g *AccountGateway) GetAllAccounts(ctx context.Context) ([]*accountResp.AccountResponse, error) {
	grpcResp, err := g.client.GetAllAccounts(
		ctx,
		mapper.ToProtoListAccount(),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToListAccountResponse(grpcResp)
	return resp, nil
}

func (g *AccountGateway) ActionOnAccount(ctx context.Context, input *accountReq.ActionRequest) (*accountResp.MessageResponse, error) {
	grpcResp, err := g.client.ActionOnAccount(
		ctx,
		mapper.ToProtoActionRequest(input),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToActionResponse(grpcResp)
	return &resp, nil
}

func (g *AccountGateway) DeleteAccount(ctx context.Context, accID string) (*accountResp.MessageResponse, error) {
	grpcResp, err := g.client.Delete(
		ctx,
		mapper.ToProtoDeleteRequest(accID),
	)

	if err != nil {
		return nil, err
	}

	resp := mapper.ToActionResponse(grpcResp)
	return &resp, nil
}
