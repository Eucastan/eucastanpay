package service

import (
	"context"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/domain"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/dto/response"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func (u *TransferUseCase) Transfer(ctx context.Context, userID string, idemKey string, input *request.TransactionIdentity) (*response.TransferResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "TransferUseCase.Transfer")
	defer span.End()

	logger := u.log.WithFields(logrus.Fields{
		"reference":       "",
		"user_id":         userID,
		"idempotency_key": idemKey,
		"amount":          input.Amount,
	})

	existing, err := u.TX.FindByIdempotencyKey(ctx, idemKey)
	switch err {

	case nil:
		resp := response.ToTransferResponse(existing)
		return &resp, nil
	case errmessage.ErrTranferNotFound:
		// continue creating transfer
	default:
		return nil, err
	}

	switch domain.TransferMode(input.Mode) {
	case domain.IntraBank, domain.InterBank, domain.OwnAccount:
	default:
		return nil, errmessage.ErrInvalidTransferMode
	}

	ref := uuid.NewString()

	transfer := &domain.Transfer{
		ID:             uuid.NewString(),
		UserID:         userID,
		Reference:      ref,
		Step:           domain.StepInitiated,
		FromAccID:      input.FromAccID,
		FromAccNo:      input.FromAccNo,
		ToAccID:        input.ToAccID,
		ToAccNo:        input.ToAccNo,
		Amount:         input.Amount,
		Description:    input.Description,
		IdempotencyKey: idemKey,
		Direction:      domain.TransferDir,
		Status:         domain.TransferStatusPending,
		Mode:           domain.TransferMode(input.Mode),
		CreatedAt:      time.Now(),
	}

	logger = logger.WithField("reference", ref)
	logger.Info("Starting transfer")

	err = u.TX.WithTx(ctx, func(tx pgx.Tx) error {

		if err := u.TX.Create(ctx, tx, transfer); err != nil {
			if err == errmessage.ErrDuplicateRequest {

				existing, findErr := u.TX.FindByIdempotencyKey(ctx, idemKey)
				if findErr != nil {
					return findErr
				}

				*transfer = *existing
				return nil
			}
			return err
		}

		transferInitiated := events.TransferInitiatedEvent{
			EventMetadata: events.NewRootEvent(ctx),
			UserID:        transfer.UserID,
			TransferID:    transfer.ID,
			Reference:     ref,
			FromAccID:     transfer.FromAccID,
			FromAccNo:     transfer.FromAccNo,
			ToAccID:       transfer.ToAccID,
			ToAccNo:       transfer.ToAccNo,
			Amount:        transfer.Amount,
			Timestamp:     time.Now().Unix(),
		}

		return u.TX.SaveOutboxEvent(ctx, tx,
			events.TopicTransferInitiated,
			transfer.Reference,
			transferInitiated,
		)
	})

	if err != nil {
		return nil, err
	}

	resp := response.ToTransferResponse(transfer)
	return &resp, err
}

func (u *TransferUseCase) RecoverStuckTransfers(ctx context.Context) {
	logger := u.log.WithField(
		"operation",
		"RecoverStuckTransfers",
	)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("stopping recovery worker")
			return

		case <-ticker.C:

			transfers, err := u.TX.FindStuckTransfers(
				ctx,
				2*time.Minute,
			)
			if err != nil {
				logger.WithError(err).
					Error("failed to find stuck transfers")
				continue
			}

			for _, t := range transfers {

				err := u.TX.WithTx(ctx, func(tx pgx.Tx) error {

					switch t.Step {

					case domain.StepInitiated:

						if err := u.TX.SaveOutboxEvent(
							ctx,
							tx,
							events.TopicDebitRequested,
							t.Reference,
							events.DebitRequestedEvent{
								EventMetadata: events.NewRootEvent(ctx),
								Reference:     t.Reference,
								FromAccID:     t.FromAccID,
								FromAccNo:     t.FromAccNo,
								ToAccNo:       t.ToAccNo,
								Amount:        t.Amount,
							},
						); err != nil {
							return err
						}

					case domain.StepDebited:

						if err := u.TX.SaveOutboxEvent(
							ctx,
							tx,
							events.TopicCreditRequested,
							t.Reference,
							events.CreditRequestedEvent{
								EventMetadata: events.NewRootEvent(ctx),
								Reference:     t.Reference,
								ToAccNo:       t.ToAccNo,
								Amount:        t.Amount,
							},
						); err != nil {
							return err
						}
					}

					return u.TX.IncrementRecoveryCount(
						ctx,
						tx,
						t.Reference,
					)
				})

				if err != nil {
					logger.
						WithField("reference", t.Reference).
						WithError(err).
						Error("failed to recover transfer")

					continue
				}

				logger.
					WithField("reference", t.Reference).
					Warn("recovery event scheduled")
			}
		}
	}
}
