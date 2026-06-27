package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/services/account/internal/domain"
	"github.com/Eucastan/eucastanpay/services/account/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/account/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/account/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type AccountUseCase struct {
	ACC    repository.AccountRepository
	logger *logrus.Logger
}

func NewAccountUseCase(acc repository.AccountRepository, logger *logrus.Logger) *AccountUseCase {
	return &AccountUseCase{
		ACC:    acc,
		logger: logger,
	}
}

func (u *AccountUseCase) CreateAccountTX(
	ctx context.Context,
	tx pgx.Tx,
	userID string,
	acc *request.CreateAccountRequest,
) (*response.AccountResponse, error) {
	u.logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"account_type": acc.AccountType,
		"service":      "account",
	}).Info("creating account")

	exists, err := u.ACC.Exists(ctx, userID, acc.AccountType)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errmessage.ErrAccAlreadyExists
	}

	account := &domain.Account{
		ID:          uuid.NewString(),
		UserID:      userID,
		Balance:     1000,
		AccountType: domain.ACCType(acc.AccountType),
		Currency:    acc.Currency,
		Status:      domain.ActiveAccount,
		CreatedAt:   time.Now(),
	}

	if err := u.ACC.Create(ctx, tx, account); err != nil {
		u.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":      userID,
			"account_type": acc.AccountType,
		}).Error("failed to create account")
		return nil, err
	}

	u.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"account_id": account.ID,
		"account_no": account.AccountNo,
	}).Info("account created")

	resp := response.ToAccountResponse(account)
	return resp, nil
}

func (u *AccountUseCase) CreateAccount(
	ctx context.Context,
	userID string,
	acc *request.CreateAccountRequest,
) (*response.AccountResponse, error) {
	u.logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"account_type": acc.AccountType,
		"service":      "account",
	}).Info("creating account")

	var a domain.Account

	err := u.ACC.WithTx(ctx, func(tx pgx.Tx) error {
		exists, err := u.ACC.Exists(ctx, userID, acc.AccountType)
		if err != nil {
			return err
		}

		if exists {
			return errmessage.ErrAccAlreadyExists
		}

		account := &domain.Account{
			ID:          uuid.NewString(),
			UserID:      userID,
			Balance:     1000, // initial unwithdrawable balance
			AccountType: domain.ACCType(acc.AccountType),
			Currency:    acc.Currency,
			Status:      domain.ActiveAccount,
			CreatedAt:   time.Now(),
		}

		if err := u.ACC.Create(ctx, tx, account); err != nil {
			u.logger.WithError(err).WithFields(logrus.Fields{
				"user_id":      userID,
				"account_id":   account.UserID,
				"account_type": account.AccountType,
			}).Error("failed to create account")
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
			u.logger.WithError(err).WithFields(logrus.Fields{
				"correlation_id": createAccountEvent.CorrelationID,
				"service":        createAccountEvent.CausationID,
				"user_id":        userID,
				"account_id":     account.UserID,
				"account_type":   account.AccountType,
			}).Error("failed to create account event outbox")

			return err
		}

		a = *account

		return nil
	})

	if err != nil {
		u.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":      userID,
			"account_id":   a.UserID,
			"account_type": a.AccountType,
		}).Error("failed account transaction")

		return nil, err
	}

	u.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"account_id": a.ID,
		"account_no": a.AccountNo,
	}).Info("account created")

	resp := response.ToAccountResponse(&a)
	return resp, nil
}

func (u *AccountUseCase) Credit(ctx context.Context, tx pgx.Tx, accID string, input *request.CreditRequest) error {
	u.logger.WithFields(logrus.Fields{
		"account_id": accID,
		"account_no": input.AccountNo,
	}).Info("initiating credit operation")

	acc, err := u.ACC.LockAccount(ctx, tx, accID, input.AccountNo)
	if err != nil {
		u.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":    acc.UserID,
			"account_id": accID,
		}).Error("failed credit operation locking")

		return err
	}

	if acc.Status == domain.CloseAccount || acc.Status == domain.FreezeAccount {
		return errmessage.ErrNotEligibleForOperation
	}

	if err := u.ACC.UpdateBalance(ctx, tx, accID, input.Amount, true); err != nil {
		return err
	}

	u.logger.WithFields(logrus.Fields{
		"account_no": input.AccountNo,
		"amount":     input.Amount,
	}).Info("account credited successfully")

	return nil
}

func (u *AccountUseCase) Debit(ctx context.Context, tx pgx.Tx, accID string, input *request.DebitRequest) error {
	u.logger.WithFields(logrus.Fields{
		"account_id": accID,
		"account_no": input.AccountNo,
	}).Info("initiating debit operation")

	acc, err := u.ACC.LockAccount(ctx, tx, accID, input.AccountNo)
	if err != nil {
		u.logger.WithError(err).WithFields(logrus.Fields{
			"account_id": accID,
			"account_no": input.AccountNo,
		}).Error("failed debit operation locking")

		return err
	}

	if acc.Status == domain.CloseAccount || acc.Status == domain.FreezeAccount {
		return errmessage.ErrNotEligibleForOperation
	}

	if acc.Balance < input.Amount {
		return errmessage.ErrInsufficientAmount
	}

	if err := u.ACC.UpdateBalance(ctx, tx, accID, input.Amount, false); err != nil {
		return err
	}

	u.logger.WithFields(logrus.Fields{
		"account_id": accID,
		"account_no": input.AccountNo,
		"amount":     input.Amount,
	}).Info("account debited successfully")

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
	acc, err := u.ACC.FindByAccountIDAndUserID(ctx, accID, userID)
	if err != nil {
		return nil, err
	}

	resp := response.ToAccountResponse(acc)

	return resp, nil
}

func (u *AccountUseCase) GetBalance(ctx context.Context, accID, userID string) (*response.AccountResponse, error) {

	acc, err := u.ACC.FindByAccountIDAndUserID(ctx, accID, userID)
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
			message = fmt.Sprintf("This account %s has been frozen", acc.ID)

		case "closed":
			acc.Status = domain.CloseAccount
			if err := u.ACC.UpdateStatus(ctx, tx, accID, string(domain.CloseAccount)); err != nil {
				return err
			}
			message = fmt.Sprintf("This account %s has been closed", acc.ID)

		default:
			return fmt.Errorf("invalid status: %s. allowed: freeze, closed", status)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return message, nil
}
