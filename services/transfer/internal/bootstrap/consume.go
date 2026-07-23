package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/consumer"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/eventhandler"
)

func (a *App) initConsumer() {
	c := consumer.NewConsumer(
		a.cfg.SharedCfg.Kafka.Brokers, a.cfg.SharedCfg.Kafka.Username,
		a.cfg.SharedCfg.Kafka.Password, "transfer-group",
		a.telemetry, a.logger,
	)

	idempotencyStore := idempotency.NewPostgresStore()

	transferConsumer := eventhandler.NewTransferConsumer(
		a.repo, idempotencyStore,
		a.telemetry, a.publish, a.logger,
	)

	c.Register(
		events.TopicTransferInitiated,
		consumer.RetryHandler(
			transferConsumer.OnTransferInitiated,
			a.publish,
			events.TopicTransferInitiated,
			events.TopicTransferDLQ,
			a.telemetry,
			3,
		),
	)

	c.Register(
		events.TopicReverseInitiated,
		consumer.RetryHandler(
			transferConsumer.OnReverseInitiated,
			a.publish,
			events.TopicReverseInitiated,
			events.TopicTransferDLQ,
			a.telemetry,
			3,
		),
	)

	c.Register(
		events.TopicDebitCompleted,
		consumer.RetryHandler(
			transferConsumer.OnDebitCompleted,
			a.publish,
			events.TopicDebitCompleted,
			events.TopicTransferDLQ,
			a.telemetry,
			3,
		),
	)

	c.Register(
		events.TopicDebitFailed,
		consumer.RetryHandler(
			transferConsumer.OnDebitFailed,
			a.publish,
			events.TopicDebitFailed,
			events.TopicTransferDLQ,
			a.telemetry,
			3,
		),
	)

	c.Register(
		events.TopicCreditCompleted,
		consumer.RetryHandler(
			transferConsumer.OnCreditCompleted,
			a.publish,
			events.TopicCreditCompleted,
			events.TopicTransferDLQ,
			a.telemetry,
			3,
		),
	)

	c.Register(
		events.TopicCreditFailed,
		consumer.RetryHandler(
			transferConsumer.OnCreditFailed,
			a.publish,
			events.TopicCreditFailed,
			events.TopicTransferDLQ,
			a.telemetry,
			3,
		),
	)

	a.consumer = c
}
