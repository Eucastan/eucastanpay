package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/clients"
	"github.com/Eucastan/eucastanpay/common/proto/account"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/domain"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type LedgerUseCase struct {
	ledger repository.LedgerRepository
	*clients.Clients
	logger *logrus.Logger
}

func NewLedgerUseCase(ledger repository.LedgerRepository, clients *clients.Clients, logger *logrus.Logger) *LedgerUseCase {
	return &LedgerUseCase{
		ledger:  ledger,
		Clients: clients,
		logger:  logger,
	}
}

func (u *LedgerUseCase) TransactionEntry(
	ctx context.Context,
	tx pgx.Tx,
	fromAccID string,
	toAccID string,
	amount int64,
	reference string,
	fromBalAfter int64,
	toBalAfter int64,
) error {

	// Debit Entry
	debit := &domain.Ledger{
		ID:           uuid.NewString(),
		AccountID:    fromAccID,
		Amount:       amount,
		EntryType:    domain.DebitEntry,
		Reference:    reference,
		BalanceAfter: fromBalAfter,
		Description:  fmt.Sprintf("Debit for transfer %s", reference),
		CreatedAt:    time.Now(),
	}

	if err := u.ledger.CreateLedgerEntry(ctx, tx, debit); err != nil {
		return err
	}

	// Credit Entry
	credit := &domain.Ledger{
		ID:           uuid.NewString(),
		AccountID:    toAccID,
		Amount:       amount,
		EntryType:    domain.CreditEntry,
		Reference:    reference,
		BalanceAfter: toBalAfter,
		Description:  fmt.Sprintf("Credit for transfer %s", reference),
		CreatedAt:    time.Now(),
	}

	if err := u.ledger.CreateLedgerEntry(ctx, tx, credit); err != nil {
		return err
	}

	// Publish events to outbox
	if err := u.publishLedgerEvent(ctx, tx, debit); err != nil {
		return err
	}
	if err := u.publishLedgerEvent(ctx, tx, credit); err != nil {
		return err
	}

	return nil
}

func (u *LedgerUseCase) publishLedgerEvent(ctx context.Context, tx pgx.Tx, entry *domain.Ledger) error {
	event := events.LedgerCreatedEvent{
		BaseEvent:    events.NewBaseEvent(ctx, "ledger-service"),
		LedgerID:     entry.ID,
		Reference:    entry.Reference,
		UserID:       "",
		AccountID:    entry.AccountID,
		Type:         string(entry.EntryType),
		Amount:       entry.Amount,
		Currency:     "NGN",
		BalanceAfter: entry.BalanceAfter,
		Description:  entry.Description,
		Timestamp:    entry.CreatedAt.Unix(),
	}

	return u.ledger.SaveOutboxEvent(ctx, tx, events.TopicLedgerCreated, entry.Reference, event)
}

func (u *LedgerUseCase) ReconcileAccount(ctx context.Context, accountID string) (*request.ReconciliationResult, error) {
	result := &request.ReconciliationResult{
		AccountID:    accountID,
		ReconciledAt: time.Now(),
		Status:       "ok",
	}

	// Get balance from Account Service
	accResp, err := u.Account.GetBalance(ctx, &account.GetBalanceRequest{Id: accountID})
	if err != nil {
		return nil, err
	}
	result.AccountBalance = accResp.Balance

	// Get ledger sum
	ledgerSum, err := u.ledger.SumByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	result.LedgerBalance = ledgerSum
	result.Difference = result.AccountBalance - result.LedgerBalance

	if result.Difference != 0 {
		result.Status = "discrepancy"
		result.Reason = "balance_mismatch_detected"

		// Publish alert event via outbox
		alertEvent := events.LedgerReconciliationAlertEvent{
			BaseEvent:      events.NewBaseEvent(ctx, "ledger-service"),
			AccountID:      accountID,
			AccountBalance: result.AccountBalance,
			LedgerBalance:  result.LedgerBalance,
			Difference:     result.Difference,
		}

		if err := u.ledger.WithTx(ctx, func(tx pgx.Tx) error {
			return u.ledger.SaveOutboxEvent(ctx, tx, events.TopicLedgerReconciliationAlert, accountID, alertEvent)
		}); err != nil {
			u.logger.WithError(err).Error("failed to save reconciliation event alert")
			return nil, err
		}
	}

	return result, nil
}

func (u *LedgerUseCase) GetTransactionEntry(ctx context.Context, id string) (*response.LedgerResponse, error) {
	ledger, err := u.ledger.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return response.ToLedgerResponse(ledger), nil
}

func (u *LedgerUseCase) GetAllLedgers(ctx context.Context) ([]response.LedgerResponse, error) {
	ledgers, err := u.ledger.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]response.LedgerResponse, 0, len(ledgers))
	for _, ledger := range ledgers {
		resp = append(resp, *response.ToLedgerResponse(&ledger))
	}

	return resp, err
}

func (u *LedgerUseCase) GetTransactionByEntryType(ctx context.Context, input *request.EntryTypeRequest) ([]response.LedgerResponse, error) {
	ledgers, err := u.ledger.FindByEntryType(ctx, input.EntryType)
	if err != nil {
		return nil, err
	}
	return response.ToListLedgerResponse(ledgers), nil
}

func (u *LedgerUseCase) GetAccountBalance(ctx context.Context, accID string) (int64, error) {
	return u.ledger.SumByAccountID(ctx, accID)
}
