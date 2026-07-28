package gateway

import (
	"context"

	ledgerpb "github.com/Eucastan/eucastanpay/common/proto/ledger"
	ledgerReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/ledger"
	ledgerResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/ledger"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/mapper"
	"github.com/sirupsen/logrus"
)

type LedgerGateway struct {
	client ledgerpb.LedgerServiceClient
	logger *logrus.Logger
}

func NewLedgerGateway(client ledgerpb.LedgerServiceClient, logger *logrus.Logger) *LedgerGateway {
	return &LedgerGateway{
		client: client,
		logger: logger,
	}
}

func (g *LedgerGateway) GetAllLedgers(ctx context.Context) ([]*ledgerResp.LedgerResponse, error) {
	grpcResp, err := g.client.GetAllLedgers(
		ctx,
		mapper.ToProtoListLedgers(),
	)

	if err != nil {
		g.logger.WithError(err).Error("failed to get ledgers: ", err.Error())
		return nil, err
	}

	resp := mapper.ToListLedgerResponse(grpcResp)
	return resp, nil
}

func (g *LedgerGateway) GetLedgerUserID(ctx context.Context, userID string) (*ledgerResp.LedgerResponse, error) {
	grpcResp, err := g.client.GetLedgerByUserId(
		ctx,
		mapper.ToProtoLedgerByUserID(userID),
	)

	if err != nil {
		g.logger.WithError(err).Error("failed to get ledgers: ", err.Error())
		return nil, err
	}

	resp := mapper.ToLedgerResponse(grpcResp)
	return resp, nil
}

func (g *LedgerGateway) GetLedger(ctx context.Context, ledgerID string) (*ledgerResp.LedgerResponse, error) {
	grpcResp, err := g.client.GetLedger(
		ctx,
		mapper.ToProtoLedger(ledgerID),
	)

	if err != nil {
		g.logger.WithError(err).Error("failed to get ledger: ", err.Error())
		return nil, err
	}

	resp := mapper.ToLedgerResponse(grpcResp)
	return resp, nil
}

func (g *LedgerGateway) GetLedgerBalance(ctx context.Context, accID string) (*ledgerResp.AccountBalanceResponse, error) {
	grpcResp, err := g.client.GetLedgerBalance(
		ctx,
		mapper.ToProtoLedgerBalance(accID),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToLedgerBalanceResponse(accID, grpcResp)
	return resp, nil
}

func (g *LedgerGateway) GetLedgerByAccountID(ctx context.Context, accID string) (*ledgerResp.LedgerResponse, error) {
	grpcResp, err := g.client.GetLedgerByAccountId(
		ctx,
		mapper.ToProtoLedgerByAccountID(accID),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToLedgerResponse(grpcResp)
	return resp, nil
}

func (g *LedgerGateway) GetLedgersByEntryType(ctx context.Context, input *ledgerReq.EntryTypeRequest) ([]*ledgerResp.LedgerResponse, error) {
	grpcResp, err := g.client.GetLedgersByEntryType(
		ctx,
		mapper.ToProtoLedgerEntryType(input),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToListLedgerEntryTypeResponse(grpcResp)
	return resp, nil
}

func (g *LedgerGateway) GetLedgerReconciliation(ctx context.Context, accID string) (*ledgerResp.ReconciliationResultResponse, error) {
	grpcResp, err := g.client.ReconcileAccount(
		ctx,
		mapper.ToProtoReconciliationByAccountID(accID),
	)

	if err != nil {
		g.logger.WithError(err).Error(err.Error())
		return nil, err
	}

	resp := mapper.ToReconciliationResponse(grpcResp)
	return resp, nil
}
