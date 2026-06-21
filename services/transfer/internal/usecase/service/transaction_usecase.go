package service

import (
	"context"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/errors"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/clients"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/proto/account"
	"github.com/Eucastan/eucastanpay/common/proto/ledger"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/domain"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type TransferUseCase struct {
	TX repository.TransferRepository
	*clients.Clients
	Publisher *producer.Publisher
	log       *logrus.Logger
}

func NewTransferUseCase(tx repository.TransferRepository, cl *clients.Clients, publisher *producer.Publisher, log *logrus.Logger) *TransferUseCase {
	return &TransferUseCase{
		TX:        tx,
		Clients:   cl,
		Publisher: publisher,
		log:       log,
	}
}

func (u *TransferUseCase) GetAllTransfers(ctx context.Context) ([]response.TransferResponse, error) {
	acc, err := u.TX.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var resp []response.TransferResponse
	for _, v := range acc {
		resp = append(resp, response.TransferResponse{
			ID:          v.ID,
			Reference:   v.Reference,
			FromAccID:   v.FromAccID,
			FromAccNo:   v.FromAccNo,
			ToAccID:     v.ToAccID,
			ToAccNo:     v.ToAccNo,
			Amount:      v.Amount,
			Description: v.Description,
			Type:        string(v.Type),
			Status:      string(v.Status),
			CreatedAt:   v.CreatedAt,
		})
	}

	return resp, err
}

func (u *TransferUseCase) GetByID(ctx context.Context, id string) (*response.TransferResponse, error) {
	logger := u.log.WithFields(logrus.Fields{
		"id": id,
	})

	transfer, err := u.TX.FindByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to transaction")
		return nil, err
	}

	resp := response.ToTransferResponse(transfer)
	return &resp, nil
}

func (u *TransferUseCase) TransferFromUser(
	ctx context.Context,
	userID string,
	idemKey string,
	input *request.TransferRequest,
) (*response.TransferResponse, error) {
	account, err := u.Account.GetUserAccount(ctx, &account.GetUserAccountRequest{
		AccountId: input.FromAccID,
		UserId:    userID,
	})
	if err != nil {
		u.log.WithError(err).Error("failed to confirm user has account")
		return nil, err
	}

	if account.UserId != userID {
		return nil, errors.ErrUserNotOwner
	}

	return u.Transfer(ctx, userID, idemKey, input)

}

func (u *TransferUseCase) ReverseTransfer(ctx context.Context, userID, originalRef, idemKey string) (*response.TransferResponse, error) {
	logger := u.log.WithFields(logrus.Fields{
		"operation":       "reverse_transfer",
		"original_ref":    originalRef,
		"idempotency_key": idemKey,
		"user_id":         userID,
	})

	existing, err := u.TX.FindByIdempotencyKey(ctx, idemKey)
	if err == nil {
		resp := response.ToTransferResponse(existing)
		return &resp, nil
	}

	var reversal *domain.Transfer

	err = u.TX.WithTx(ctx, func(tx pgx.Tx) error {
		// Find original transfer inside transaction
		original, err := u.TX.FindByReference(ctx, tx, originalRef)
		if err != nil {
			return err
		}

		if original.IsReversed {
			return errors.ErrAlreadyReversed
		}
		if original.Status != domain.TransferStatusSuccess {
			return errors.ErrCannotReverseNonSuccessfulTransfer
		}

		// Create reversal transfer
		reversalRef := uuid.NewString()
		reversal = &domain.Transfer{
			ID:             uuid.NewString(),
			UserID:         userID,
			Reference:      reversalRef,
			Step:           domain.StepInitiated,
			FromAccID:      original.ToAccID,
			FromAccNo:      original.ToAccNo,
			ToAccID:        original.FromAccID,
			ToAccNo:        original.FromAccNo,
			Amount:         original.Amount,
			Description:    "Reversal of " + originalRef,
			IdempotencyKey: idemKey,
			Type:           domain.TransferTypeReverse,
			Status:         domain.TransferStatusPending,
			Mode:           original.Mode,
			ReversalRef:    originalRef,
			CreatedAt:      time.Now(),
		}

		if err := u.TX.Create(ctx, tx, reversal); err != nil {
			return err
		}

		// Mark original transfer as reversing
		if err := u.TX.UpdateStatus(ctx, tx, originalRef, string(domain.TransferStatusReversing)); err != nil {
			return err
		}

		// Publish reversal event to start saga
		reverseEvent := events.ReverseDebitEvent{
			BaseEvent: events.NewBaseEvent(ctx, "transfer-service"),
			Reference: reversalRef,
			AccountID: reversal.FromAccID,
			AccountNo: reversal.FromAccNo,
			Amount:    reversal.Amount,
		}

		return u.TX.SaveOutboxEvent(ctx, tx, events.TopicDebitReverse, reversalRef, reverseEvent)
	})

	if err != nil {
		logger.WithError(err).Error("Reverse transfer failed")
		return nil, err
	}

	logger.Info("Reversal initiated successfully")
	resp := response.ToTransferResponse(reversal)
	return &resp, nil
}

func (u *TransferUseCase) ReconcileAccount(
	ctx context.Context,
	accID string,
	accNo int64,
) error {

	logger := u.log.WithField("account_id", accID)

	// Get balance from Account Service
	_, err := u.Account.GetBalance(ctx, &account.GetBalanceRequest{
		Id:        accID,
		AccountNo: accNo,
	})

	if err != nil {
		logger.WithError(err).Error("failed to confirm account exists during reconciliation")
		return err
	}

	// Get balance from Ledger Service
	_, err = u.Ledger.ReconcileAccount(ctx, &ledger.ReconcileAccountRequest{AccountId: accID})
	if err != nil {
		logger.WithError(err).Error("failed to fetch ledger balance")
		return err
	}

	if err := u.TX.WithTx(ctx, func(tx pgx.Tx) error {
		auditEvent := events.AdminActionEvent{
			BaseEvent:  events.NewBaseEvent(ctx, "transfer-service"),
			AdminID:    "system", // or extract from context
			Action:     "reconcile_account",
			TargetType: "account",
			TargetID:   accID,
			Reason:     "manual_reconciliation",
		}
		return u.TX.SaveOutboxEvent(ctx, tx, events.TopicAdminActionTaken, accID, auditEvent)
	}); err != nil {
		logger.WithError(err).Error("failed to save admin action event")
	}

	logger.Info("reconciliation successful")

	return nil
}
