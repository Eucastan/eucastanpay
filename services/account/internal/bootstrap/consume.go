package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/consumer"
	"github.com/Eucastan/eucastanpay/services/account/internal/eventhandler"
)

func (a *App) initConsumer() {

	c := consumer.NewConsumer(
		a.cfg.SharedCfg.Kafka.Brokers, a.cfg.SharedCfg.Kafka.Username,
		a.cfg.SharedCfg.Kafka.Password, "account-service-group",
		a.telemetry, a.logger,
	)

	idempotencyStore := idempotency.NewPostgresStore()
	accountConsumer := eventhandler.NewAccountConsumer(
		a.repo, a.uc, idempotencyStore,
		a.publish, a.telemetry, a.logger,
	)

	c.Register(events.TopicUserRegistered,
		consumer.RetryHandler(
			accountConsumer.OnCreateAccountRequest,
			a.publish,
			events.TopicUserRegistered,
			events.TopicAccountDLQ,
			a.telemetry,
			3,
		),
	)

	c.Register(events.TopicDebitRequested,
		consumer.RetryHandler(
			accountConsumer.OnDebitRequested,
			a.publish,
			events.TopicDebitRequested,
			events.TopicAccountDLQ,
			a.telemetry,
			3,
		),
	)

	c.Register(events.TopicCreditRequested,
		consumer.RetryHandler(
			accountConsumer.OnCreditRequested,
			a.publish,
			events.TopicCreditRequested,
			events.TopicAccountDLQ,
			a.telemetry,
			3,
		),
	)

	a.consumer = c
}
