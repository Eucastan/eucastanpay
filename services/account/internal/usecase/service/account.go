package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/errors"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/services/account/internal/domain"
	"github.com/Eucastan/eucastanpay/services/account/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/account/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/account/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AccountUseCase struct {
	ACC       repository.AccountRepository
	Publisher *producer.Publisher
}

func NewAccountUseCase(acc repository.AccountRepository, publisher *producer.Publisher) *AccountUseCase {
	return &AccountUseCase{
		ACC:       acc,
		Publisher: publisher,
	}
}

func (u *AccountUseCase) CreateAccount(
	ctx context.Context,
	userID string,
	acc *request.CreateAccountRequest,
) (*response.AccountResponse, error) {

	account := &domain.Account{
		ID:          uuid.NewString(),
		UserID:      userID,
		AccountNo:   acc.AccountNo,
		Balance:     acc.Balance,
		AccountType: domain.ACCType(acc.AccountType),
		Currency:    acc.Currency,
		Status:      domain.ActiveAccount,
		CreatedAt:   time.Now(),
	}

	err := u.ACC.WithTx(ctx, func(tx pgx.Tx) error {

		if err := u.ACC.Create(ctx, tx, account); err != nil {
			return err
		}

		createAccountEvent := events.AccountCreatedEvent{
			BaseEvent:   events.NewBaseEvent(ctx, "account-service"),
			AccountID:   account.ID,
			UserID:      account.UserID,
			AccountNo:   account.AccountNo,
			AccountType: string(account.AccountType),
			Currency:    account.Currency,
			Timestamp:   time.Now().Unix(),
		}

		if err := u.ACC.SaveOutboxEvent(ctx, tx, events.TopicAccountCreated, account.ID, createAccountEvent); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := response.ToAccountResponse(account)
	return resp, nil
}

func (u *AccountUseCase) Credit(ctx context.Context, tx pgx.Tx, accID string, input *request.CreditRequest) error {

	acc, err := u.ACC.LockAccount(ctx, tx, accID, input.AccountNo)
	if err != nil {
		return err
	}

	if acc.Status == domain.CloseAccount || acc.Status == domain.FreezeAccount {
		return errors.ErrNotEligibleForOperation
	}

	if err := u.ACC.UpdateBalance(ctx, tx, accID, input.Amount, true); err != nil {
		return err
	}

	return nil
}

func (u *AccountUseCase) Debit(ctx context.Context, tx pgx.Tx, accID string, input *request.DebitRequest) error {

	acc, err := u.ACC.LockAccount(ctx, tx, accID, input.AccountNo)
	if err != nil {
		return err
	}

	if acc.Status == domain.CloseAccount || acc.Status == domain.FreezeAccount {
		return errors.ErrNotEligibleForOperation
	}

	if acc.Balance < input.Amount {
		return errors.ErrInsufficientAmount
	}

	if err := u.ACC.UpdateBalance(ctx, tx, accID, input.Amount, false); err != nil {
		return err
	}

	return nil

}

func (u *AccountUseCase) GetAllAccount(ctx context.Context) ([]response.AccountResponse, error) {
	accounts, err := u.ACC.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]response.AccountResponse, 0, len(accounts))
	for _, acc := range accounts {
		resp = append(resp, *response.ToAccountResponse(&acc))
	}

	return resp, err
}

func (u *AccountUseCase) GetByUserID(ctx context.Context, userID string) (*response.AccountResponse, error) {
	acc, err := u.ACC.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := response.ToAccountResponse(acc)

	return resp, nil
}

func (u *AccountUseCase) GetByAccountIDAndUserID(ctx context.Context, accID, userID string) (*response.AccountResponse, error) {
	acc, err := u.ACC.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := response.ToAccountResponse(acc)

	return resp, nil
}

func (u *AccountUseCase) GetBalance(ctx context.Context, accID string) (*response.AccountResponse, error) {

	acc, err := u.ACC.FindByID(ctx, accID)
	if err != nil {
		return nil, err
	}

	return response.ToAccountResponse(acc), nil
}

func (u *AccountUseCase) ActionOnAccount(ctx context.Context, accID, status string, accNo int64) (string, error) {
	var message string
	err := u.ACC.WithTx(ctx, func(tx pgx.Tx) error {
		acc, err := u.ACC.FindByIDTX(ctx, tx, accID, accNo)
		if err != nil {
			return err
		}

		// Only allow status changes if the account is currently Active
		if acc.Status != domain.ActiveAccount {
			return fmt.Errorf("only active accounts can change status. current status: %s", acc.Status)
		}

		switch status {
		case "freeze":
			acc.Status = domain.FreezeAccount
			if err := u.ACC.UpdateStatus(ctx, tx, accID, string(domain.FreezeAccount)); err != nil {
				return err
			}
			message = fmt.Sprintf("This account %s has been freezed", acc.ID)

		case "closed":
			acc.Status = domain.CloseAccount
			if err := u.ACC.UpdateStatus(ctx, tx, accID, string(domain.CloseAccount)); err != nil {
				return err
			}
			message = fmt.Sprintf("This account %s has been closed", acc.ID)

		default:
			return fmt.Errorf("invalid status: %s. allowed: freezed, closed", status)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return message, nil
}
