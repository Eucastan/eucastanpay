package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/consumer"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/eventshandler"
)

func (a *App) initConsumer() {
	c := consumer.NewConsumer(
		a.cfg.SharedCfg.Kafka.Brokers,
		a.cfg.SharedCfg.Kafka.Username,
		a.cfg.SharedCfg.Kafka.Password,
		"ledger-service-group",
		a.telemetry,
		a.logger,
	)

	idempotency := idempotency.NewPostgresStore()
	ledgerConsumer := eventshandler.NewLedgerEventHandler(
		a.repo, a.uc, a.telemetry,
		idempotency, a.publish, a.logger,
	)

	c.Register(events.TopicTransferCompleted,
		consumer.RetryHandler(
			ledgerConsumer.OnTransferCompleted,
			a.publish,
			events.TopicTransferCompleted,
			events.TopicLedgerDLQ,
			a.telemetry,
			3,
		),
	)

	c.Register(events.TopicDepositAccount,
		consumer.RetryHandler(
			ledgerConsumer.OnAccountDeposit,
			a.publish,
			events.TopicDepositAccount,
			events.TopicLedgerDLQ,
			a.telemetry,
			3,
		),
	)

	c.Register(events.TopicWithdrawal,
		consumer.RetryHandler(
			ledgerConsumer.OnCasWithdraw,
			a.publish,
			events.TopicWithdrawal,
			events.TopicLedgerDLQ,
			a.telemetry,
			3,
		),
	)

	a.consumer = c
}
