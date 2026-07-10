package service

import (
	"context"
	"fmt"
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

	resp := make([]response.TransferResponse, 0, len(acc))
	for _, v := range acc {
		resp = append(resp, response.ToTransferResponse(&v))
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

	acc, err := u.Account.GetUserAccount(ctx, &account.GetUserAccountRequest{
		UserId: userID,
	})
	if err != nil {
		u.log.WithError(err).Error("failed to confirm user has account")
		u.telemetry.RecordError(span, err)
		return nil, err
	}

	confirm, err := u.Account.ResolveAccount(ctx, &account.ConfirmAccountRequest{
		FromAccountNo: acc.AccountNo,
		ToAccountNo:   input.ToAccNo,
	})

	if err != nil {
		return nil, err
	}

	if acc.UserId != userID {
		return nil, errmessage.ErrUserNotOwner
	}

	identity := &request.TransactionIdentity{
		FromAccID:   acc.AccountId,
		FromAccNo:   acc.AccountNo,
		ToAccID:     confirm.ToAccountId,
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
		desc := fmt.Sprintf("Reversal of %s from %d to %d", originalRef, original.ToAccNo, original.FromAccNo)
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
			Description:    desc,
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
		reverseEvent := events.ReverseInitiatedEvent{
			EventMetadata: events.NewRootEvent(ctx),
			UserID:        reversal.UserID,
			TransferID:    reversal.ID,
			Reference:     reversalRef,
			FromAccID:     reversal.FromAccID,
			FromAccNo:     reversal.FromAccNo,
			ToAccID:       reversal.ToAccID,
			ToAccNo:       reversal.ToAccNo,
			Amount:        reversal.Amount,
		}

		return u.TX.SaveOutboxEvent(ctx, tx, events.TopicReverseInitiated, reversalRef, reverseEvent)
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
	input *request.ReconciliationRequest,
) error {
	ctx, span := u.telemetry.Start(ctx, "TransferUseCase.ReconcileAccount")
	defer span.End()

	logger := u.log.WithField("account_id", accID)

	// Get balance from Account Service
	_, err := u.Account.GetBalance(ctx, &account.GetBalanceRequest{
		Id:        accID,
		AccountNo: input.AccountNo,
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
