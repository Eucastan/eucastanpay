package bootstrap

import "github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"

func (a *App) initPublisher() {
	a.publish = producer.NewPublisher(
		a.cfg.Kafka.Brokers,
		a.cfg.Kafka.Username,
		a.cfg.Kafka.Password,
		a.telemetry,
	)
}
