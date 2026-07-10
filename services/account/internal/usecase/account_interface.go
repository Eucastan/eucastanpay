package usecase

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/account/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/account/internal/dto/response"
	"github.com/jackc/pgx/v5"
)

type AccountUseCase interface {
	CreateAccountTX(ctx context.Context, tx pgx.Tx, userID, email string, acc *request.CreateAccountRequest) (*response.AccountResponse, error)
	CreateAccount(ctx context.Context, userID, email string, acc *request.CreateAccountRequest) (*response.AccountResponse, error)
	DepositAccount(ctx context.Context, accID string, input *request.DepositRequest) error
	Credit(ctx context.Context, tx pgx.Tx, accID string, input *request.CreditRequest) error
	Debit(ctx context.Context, tx pgx.Tx, accID string, input *request.DebitRequest) error
	WithDrawal(ctx context.Context, accID string, input *request.DepositRequest) error
	GetAllAccount(ctx context.Context) ([]response.AccountResponse, error)
	GetByUserID(ctx context.Context, userID string) (*response.AccountResponse, error)
	ConfirmSenderAndReceiver(ctx context.Context, fromAccNo int64, toAccNo int64) (*response.ConfirmAccountResponse, error)
	GetByAccountIDAndUserID(ctx context.Context, accID, userID string) (*response.AccountResponse, error)
	GetBalance(ctx context.Context, accID, userID string) (*response.AccountResponse, error)
	ActionOnAccount(ctx context.Context, accID, status string, accNo int64) (string, error)
}
