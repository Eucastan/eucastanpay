package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/consumer"
	"github.com/Eucastan/eucastanpay/services/audit/internal/eventhandler"
)

func (a *App) initConsumer() {
	c := consumer.NewConsumer(
		a.cfg.SharedCfg.Kafka.Brokers,
		a.cfg.SharedCfg.Kafka.Username,
		a.cfg.SharedCfg.Kafka.Password,
		"audit-service-group",
		a.telemetry, a.logger,
	)

	idempotencyStore := idempotency.NewPostgresStore()
	auditConsumer := eventhandler.NewAuditConsumer(
		a.repo, idempotencyStore,
		a.telemetry, a.logger,
	)

	// Register multiple topics
	topics := []string{
		events.TopicUserRegistered,
		events.TopicUserRegistrationFailed,
		events.TopicUserKYCCreated,
		events.TopicUserKYCVerified,
		events.TopicAccountCreated,
		events.TopicCreateAccFailed,
		events.TopicDepositAccount,
		events.TopicWithdrawal,
		events.TopicTransferInitiated,
		events.TopicReverseInitiated,
		events.TopicTransferCompleted,
		events.TopicTransferFailed,
		events.TopicDebitCompleted,
		events.TopicCreditCompleted,
		events.TopicLedgerCreated,
		events.TopicDebitRequested,
		events.TopicCreditRequested,
	}

	for _, topic := range topics {
		c.Register(topic, consumer.RetryHandler(
			auditConsumer.Handler(topic),
			a.publish,
			topic,
			events.TopicAuditDLQ,
			a.telemetry,
			3,
		))
	}

	a.consumer = c
}
