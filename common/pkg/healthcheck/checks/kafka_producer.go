package checks

import (
	"context"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
)

type KafkaProducer struct {
	producer healthcheck.KafkaProducer
}

func NewKafkaProducer(p healthcheck.KafkaProducer) *KafkaProducer {
	return &KafkaProducer{
		producer: p,
	}
}

func (k *KafkaProducer) Name() string {
	return "kafka-producer"
}

func (k *KafkaProducer) Check(ctx context.Context) healthcheck.Component {
	started := time.Now()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := k.producer.Ping(ctx); err != nil {
		return healthcheck.Component{
			Name:     k.Name(),
			Status:   healthcheck.Unhealthy,
			Error:    err.Error(),
			Duration: time.Since(started).String(),
		}
	}

	return healthcheck.Component{
		Name:     k.Name(),
		Status:   healthcheck.Healthy,
		Duration: time.Since(started).String(),
	}
}
