package healthcheck

import "context"

type KafkaProducer interface {
	Ping(ctx context.Context) error
}

type KafkaConsumer interface {
	Ping(ctx context.Context) error
}

type GRPCConnection interface {
	Ping(ctx context.Context) error
}
