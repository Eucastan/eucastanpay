package usecase

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/ledger/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/dto/response"
	"github.com/jackc/pgx/v5"
)

type LedgerUseCase interface {
	CreateEntries(ctx context.Context, tx pgx.Tx, userID string, fromAccID string, toAccID string, amount int64, reference string, fromBalAfter int64, toBalAfter int64) error
	ReconcileAccount(ctx context.Context, accountID string) (*response.ReconciliationResult, error)
	GetTransactionEntry(ctx context.Context, id string) (*response.LedgerResponse, error)
	GetAllLedgers(ctx context.Context) ([]response.LedgerResponse, error)
	GetTransactionByEntryType(ctx context.Context, input *request.EntryTypeRequest) ([]response.LedgerResponse, error)
	GetAccountBalance(ctx context.Context, accID string) (int64, error)
}
