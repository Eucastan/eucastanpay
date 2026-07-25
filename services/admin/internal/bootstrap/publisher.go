package bootstrap

import "github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"

func (a *App) initPublisher() {
	a.publish = producer.NewPublisher(
		a.cfg.SharedCfg.Kafka.Brokers,
		a.cfg.SharedCfg.Kafka.Username,
		a.cfg.SharedCfg.Kafka.Password,
		a.telemetry,
	)
}
