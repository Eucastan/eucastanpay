package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/Eucastan/eucastanpay/common/pkg/config"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/kafkaconfig"
	"github.com/segmentio/kafka-go"
)

type Initializer struct {
	cfg config.KafkaConfig
}

func NewInitializer(cfg config.KafkaConfig) *Initializer {
	return &Initializer{
		cfg: cfg,
	}
}

func (i *Initializer) EnsureTopics(ctx context.Context) error {
	topics := []string{
		events.TopicUserRegistered,
		events.TopicUserKYCCreated,
		events.TopicUserKYCVerified,

		events.TopicAccountCreated,
		events.TopicDepositAccount,
		events.TopicWithdrawal,

		events.TopicTransferInitiated,
		events.TopicTransferReversed,
		events.TopicReverseInitiated,
		events.TopicDebitRequested,
		events.TopicDebitCompleted,
		events.TopicCreditRequested,
		events.TopicCreditCompleted,
		events.TopicTransferCompleted,
		events.TopicTransferRetry,
		events.TopicLedgerCreated,
		events.TopicAuditCreated,

		events.TopicAdminActionTaken,
		events.TopicUserStatusChanged,
		events.TopicAccountStatusChanged,

		events.TopicUserRegistrationFailed,
		events.TopicCreateAccFailed,
		events.TopicTransferFailed,
		events.TopicCreditFailed,
		events.TopicDebitFailed,
		events.TopicAuditFailed,

		events.TopicLedgerReconciliationAlert,

		events.TopicTransferDLQ,
		events.TopicAccountDLQ,
		events.TopicAdminDLQ,
		events.TopicLedgerDLQ,
		events.TopicAuditDLQ,
		events.TopicNotificationDLQ,
	}

	dialer := kafkaconfig.NewDialer(i.cfg)

	if len(i.cfg.Brokers) == 0 {
		return fmt.Errorf("no kafka brokers configured")
	}

	conn, err := dialer.DialContext(ctx, "tcp", i.cfg.Brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	// Retrieve the controller
	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	// Connect to controller
	controllerHostAndPort := fmt.Sprintf("%s:%d", controller.Host, controller.Port)
	controllerConn, err := dialer.DialContext(ctx, "tcp", controllerHostAndPort)
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicCfg := make([]kafka.TopicConfig, len(topics))
	for i, topic := range topics {
		topicCfg[i] = kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	err = controllerConn.CreateTopics(topicCfg...)
	if err != nil &&
		!strings.Contains(strings.ToLower(err.Error()), "already exists") {
		return err
	}

	return nil
}
