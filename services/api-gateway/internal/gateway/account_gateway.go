package gateway

import (
	"context"

	accountpb "github.com/Eucastan/eucastanpay/common/proto/account"
	accountReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/account"
	accountResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/account"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/mapper"
	"github.com/sirupsen/logrus"
)

type AccountGateway struct {
	client accountpb.AccountServiceClient
	logger *logrus.Logger
}

func NewAccountGateway(client accountpb.AccountServiceClient, logger *logrus.Logger) *AccountGateway {
	return &AccountGateway{
		client: client,
		logger: logger,
	}
}

func (g *AccountGateway) Deposit(ctx context.Context, accID string, input *accountReq.DepositRequest) (*accountResp.MessageResponse, error) {
	grpcResp, err := g.client.Deposit(
		ctx,
		mapper.ToProtoDeposit(accID, *input),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToActionResponse(grpcResp)
	return &resp, nil
}

func (g *AccountGateway) WithDraw(ctx context.Context, accID string, input *accountReq.DepositRequest) (*accountResp.MessageResponse, error) {
	grpcResp, err := g.client.Withdraw(
		ctx,
		mapper.ToProtoWithdraw(accID, *input),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToActionResponse(grpcResp)
	return &resp, nil
}

func (g *AccountGateway) GetBalance(ctx context.Context, accID, userID string) (*accountResp.AccountResponse, error) {
	grpcResp, err := g.client.GetBalance(
		ctx,
		mapper.ToProtoGetBalance(accID, userID),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
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
		g.logger.WithError(err).Error(err.Error())
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
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToListAccountResponse(grpcResp)
	return resp, nil
}

func (g *AccountGateway) ActionOnAccount(ctx context.Context, accID string, input *accountReq.ActionRequest) (*accountResp.MessageResponse, error) {
	grpcResp, err := g.client.ActionOnAccount(
		ctx,
		mapper.ToProtoActionRequest(accID, input),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
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
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToActionResponse(grpcResp)
	return &resp, nil
}
