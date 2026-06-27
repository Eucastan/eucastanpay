package usecase

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/account/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/account/internal/dto/response"
	"github.com/jackc/pgx/v5"
)

type AccountUseCase interface {
	CreateAccountTX(ctx context.Context, tx pgx.Tx, userID string, acc *request.CreateAccountRequest) (*response.AccountResponse, error)
	CreateAccount(ctx context.Context, userID string, acc *request.CreateAccountRequest) (*response.AccountResponse, error)
	Credit(ctx context.Context, tx pgx.Tx, accID string, input *request.CreditRequest) error
	Debit(ctx context.Context, tx pgx.Tx, accID string, input *request.DebitRequest) error
	GetAllAccount(ctx context.Context) ([]response.AccountResponse, error)
	GetByUserID(ctx context.Context, userID string) (*response.AccountResponse, error)
	GetByAccountIDAndUserID(ctx context.Context, accID, userID string) (*response.AccountResponse, error)
	GetBalance(ctx context.Context, accID, userID string) (*response.AccountResponse, error)
	ActionOnAccount(ctx context.Context, accID, status string, accNo int64) (string, error)
}
