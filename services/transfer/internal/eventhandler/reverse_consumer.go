package eventhandler

import (
	"context"
	"fmt"

	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type SagaRequest struct {
	ParentMetadata events.EventMetadata
	UserID string
	Reference      string
	FromAccID      string
	FromAccNo      int64
	ToAccID        string
	ToAccNo        int64
	Amount         int64
	ProcessedTopic string
}

func (h *TransferConsumer) OnReverseInitiated(ctx context.Context, msg []byte) error {
	ctx, span := h.telemetry.Start(ctx, "TransferConsumer.OnReverseInitiated")
	defer span.End()

	h.logger.Info("Reverse Initiated Received")

	event, err := kafka.Decode[events.ReverseInitiatedEvent](msg)
	if err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"operation": "OnReverseInitiated",
		"reference": event.Reference,
	}).Info("Entering the transactional block")

	return h.emitDebitRequest(ctx, SagaRequest{
		ParentMetadata: event.EventMetadata,
		UserID: event.UserID,
		Reference:      event.Reference,
		FromAccID:      event.FromAccID,
		ToAccID:        event.ToAccID,
		FromAccNo:      event.FromAccNo,
		ToAccNo:        event.ToAccNo,
		Amount:         event.Amount,
		ProcessedTopic: events.TopicReverseInitiated,
	})
}

func (h *TransferConsumer) emitDebitRequest(
	ctx context.Context,
	req SagaRequest,
) error {
	ctx, span := h.telemetry.Start(ctx, "TransferConsumer.startSaga")
	defer span.End()

	eventID := fmt.Sprintf("%s:%s", req.Reference, req.ProcessedTopic)
	return h.repo.WithTx(ctx, func(tx pgx.Tx) error {
		processed, err := h.idemStore.IsEventProcessedTx(ctx, tx, eventID)
		if err != nil {
			span.RecordError(err)
			return err
		}

		if processed {
			return events.ErrProcessed
		}

		err = h.repo.SaveOutboxEvent(ctx, tx, events.TopicDebitRequested, req.Reference,
			events.DebitRequestedEvent{
				EventMetadata: events.NewChildEvent(req.ParentMetadata),
				UserID: req.UserID,
				Reference:     req.Reference,
				FromAccID:     req.FromAccID,
				FromAccNo:     req.FromAccNo,
				ToAccID:       req.ToAccID,
				ToAccNo:       req.ToAccNo,
				Amount:        req.Amount,
			},
		)
		if err != nil {
			span.RecordError(err)
			return err
		}

		return h.idemStore.MarkEventProcessedTx(
			ctx, tx, uuid.NewString(),
			eventID, req.ProcessedTopic,
		)
	})
}
