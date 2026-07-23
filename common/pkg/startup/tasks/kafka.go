package tasks

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
)

type Kafka struct {
	producer healthcheck.KafkaProducer
}

func (k *Kafka) Run(ctx context.Context) error {
	return k.producer.Ping(ctx)
}
