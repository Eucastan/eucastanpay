package eventshandler

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/repository"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/usecase"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type LedgerEventHandler struct {
	Repo      repository.LedgerRepository
	Ledger    usecase.LedgerUseCase
	telemetry *telemetry.Telemetry
	IdemStore idempotency.Store
	Publisher *producer.Publisher
	Log       *logrus.Logger
}

func NewLedgerEventHandler(
	repo repository.LedgerRepository,
	ledger usecase.LedgerUseCase,
	telemetry *telemetry.Telemetry,
	IdemStore idempotency.Store,
	publisher *producer.Publisher,
	log *logrus.Logger) *LedgerEventHandler {

	return &LedgerEventHandler{
		Repo:      repo,
		Ledger:    ledger,
		telemetry: telemetry,
		IdemStore: IdemStore,
		Publisher: publisher,
		Log:       log,
	}
}

func (h *LedgerEventHandler) OnTransferCompleted(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "LedgerEventHandler.OnTransferCompleted")
	defer span.End()

	event, err := kafka.Decode[events.TransferCompletedEvent](msg)
	if err != nil {
		span.RecordError(err)
		return err
	}

	eventID := event.Reference + ":" + events.TopicLedgerCreated
	return h.Repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.IdemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			return err
		}
		if processed {
			return nil
		}

		err = h.Ledger.TransactionEntry(
			ctx,
			tx,
			event.FromAccID,
			event.ToAccID,
			event.Amount,
			event.Reference,
			event.FromBalanceAfter,
			event.ToBalanceAfter,
		)
		if err != nil {
			span.RecordError(err)
			return err
		}

		return h.IdemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), eventID, events.TopicTransferCompleted)
	})

}
