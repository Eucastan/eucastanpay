package service

import (
	"context"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/clients"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
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
	telemetry *telemetry.Telemetry
	log       *logrus.Logger
}

func NewTransferUseCase(tx repository.TransferRepository, cl *clients.Clients, publisher *producer.Publisher, telemetry *telemetry.Telemetry, log *logrus.Logger) *TransferUseCase {
	return &TransferUseCase{
		TX:        tx,
		Clients:   cl,
		Publisher: publisher,
		telemetry: telemetry,
		log:       log,
	}
}

func (u *TransferUseCase) GetAllTransfers(ctx context.Context) ([]response.TransferResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "TransferUseCase.GetAllTransfers")
	defer span.End()

	acc, err := u.TX.FindAll(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	var resp []response.TransferResponse
	for _, v := range acc {
		resp = append(resp, response.TransferResponse{
			ID:          v.ID,
			Reference:   v.Reference,
			FromAccID:   v.FromAccID,
			FromAccNo:   v.FromAccNo,
			ToAccNo:     v.ToAccNo,
			Amount:      v.Amount,
			Description: v.Description,
			Direction:   string(v.Direction),
			Status:      string(v.Status),
			CreatedAt:   v.CreatedAt,
		})
	}

	return resp, err
}

func (u *TransferUseCase) GetByID(ctx context.Context, id string) (*response.TransferResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "TransferUseCase.GetByID")
	defer span.End()

	logger := u.log.WithFields(logrus.Fields{
		"operation": "GetByID",
		"id":        id,
	})

	transfer, err := u.TX.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
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
	ctx, span := u.telemetry.Start(ctx, "TransferUseCase.TransferFromUser")
	defer span.End()

	u.telemetry.Count(ctx, 1)

	account, err := u.Account.GetUserAccount(ctx, &account.GetUserAccountRequest{
		UserId: userID,
	})
	if err != nil {
		u.log.WithError(err).Error("failed to confirm user has account")
		u.telemetry.RecordError(span, err)
		return nil, err
	}

	if account.UserId != userID {
		return nil, errmessage.ErrUserNotOwner
	}

	identity := &request.TransactionIdentity{
		FromAccID:   account.AccountId,
		FromAccNo:   account.AccountNo,
		ToAccNo:     input.ToAccNo,
		Amount:      input.Amount,
		Description: input.Description,
		Mode:        input.Mode,
	}

	return u.Transfer(ctx, userID, idemKey, identity)

}

func (u *TransferUseCase) ReverseTransfer(ctx context.Context, userID, originalRef, idemKey string) (*response.TransferResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "TransferUseCase.ReverseTransfer")
	defer span.End()

	logger := u.log.WithFields(logrus.Fields{
		"operation":       "reverse_transfer",
		"original_ref":    originalRef,
		"idempotency_key": idemKey,
		"user_id":         userID,
	})

	existing, err := u.TX.FindByIdempotencyKey(ctx, idemKey)
	if err == nil {
		span.RecordError(err)
		resp := response.ToTransferResponse(existing)
		return &resp, nil
	}

	var reversal *domain.Transfer

	err = u.TX.WithTx(ctx, func(tx pgx.Tx) error {
		// Find original transfer inside transaction
		original, err := u.TX.FindByReference(ctx, tx, originalRef)
		if err != nil {
			span.RecordError(err)
			return err
		}

		if original.IsReversed {
			return errmessage.ErrAlreadyReversed
		}
		if original.Status != domain.TransferStatusSuccess {
			return errmessage.ErrCannotReverseNonSuccessfulTransfer
		}

		// Create reversal transfer
		reversalRef := uuid.NewString()
		reversal = &domain.Transfer{
			ID:             uuid.NewString(),
			UserID:         userID,
			Reference:      reversalRef,
			Step:           domain.StepInitiated,
			FromAccNo:      original.ToAccNo,
			ToAccNo:        original.FromAccNo,
			Amount:         original.Amount,
			Description:    "Reversal of " + originalRef,
			IdempotencyKey: idemKey,
			Direction:      domain.ReverseDir,
			Status:         domain.TransferStatusPending,
			Mode:           original.Mode,
			ReversalRef:    originalRef,
			CreatedAt:      time.Now(),
		}

		if err := u.TX.Create(ctx, tx, reversal); err != nil {
			span.RecordError(err)
			return err
		}

		// Mark original transfer as reversing
		if err := u.TX.UpdateStatus(ctx, tx, originalRef, string(domain.TransferStatusReversing)); err != nil {
			span.RecordError(err)
			return err
		}

		// Publish reversal event to start saga
		reverseEvent := events.ReverseDebitEvent{
			EventMetadata: events.NewRootEvent(ctx),
			Reference:     reversalRef,
			AccountID:     reversal.FromAccID,
			AccountNo:     reversal.FromAccNo,
			Amount:        reversal.Amount,
		}

		return u.TX.SaveOutboxEvent(ctx, tx, events.TopicDebitReverse, reversalRef, reverseEvent)
	})

	if err != nil {
		span.RecordError(err)
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
	ctx, span := u.telemetry.Start(ctx, "TransferUseCase.ReconciliationAccount")
	defer span.End()

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

	err = u.TX.WithTx(ctx, func(tx pgx.Tx) error {
		auditEvent := events.AdminActionEvent{
			EventMetadata: events.NewRootEvent(ctx),
			AdminID:       "system",
			Action:        "reconcile_account",
			TargetType:    "account",
			TargetID:      accID,
			Reason:        "manual_reconciliation",
		}
		return u.TX.SaveOutboxEvent(ctx, tx, events.TopicAdminActionTaken, accID, auditEvent)
	})
	if err != nil {
		logger.WithError(err).Error("failed to save admin action event")
		return err
	}

	logger.Info("reconciliation successful")

	return nil
}
