package eventshandler

import (
	"context"
	"encoding/json"

	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/repository"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/usecase"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type LedgerEventHandler struct {
	Repo      repository.LedgerRepository
	Ledger    usecase.LedgerUseCase
	IdemStore idempotency.Store
	Publisher *producer.Publisher
	Log       *logrus.Logger
}

func NewLedgerEventHandler(
	repo repository.LedgerRepository,
	ledger usecase.LedgerUseCase,
	IdemStore idempotency.Store,
	publisher *producer.Publisher,
	log *logrus.Logger) *LedgerEventHandler {

	return &LedgerEventHandler{
		Repo:      repo,
		Ledger:    ledger,
		IdemStore: IdemStore,
		Publisher: publisher,
		Log:       log,
	}
}

func (h *LedgerEventHandler) OnTransferCompleted(ctx context.Context, msg []byte) error {
	var event events.TransferCompletedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	eventID := event.Reference + ":" + events.TopicLedgerCreated
	return h.Repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.IdemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
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
			return err
		}

		return h.IdemStore.MarkEventProcessedTx(ctx, tx, uuid.NewString(), eventID, events.TopicTransferCompleted)
	})

}
