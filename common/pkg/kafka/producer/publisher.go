package producer

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/sirupsen/logrus"

	"github.com/Eucastan/eucastanpay/common/pkg/config"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/kafkaconfig"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Publisher struct {
	writer    *kafka.Writer
	brokers   []string
	telemetry *telemetry.Telemetry
	mechanism *plain.Mechanism
	tlsConfig *tls.Config
	logger    *logrus.Logger
	TLS       bool
	SASL      bool
}

func NewPublisher(cfg config.KafkaConfig, telemetry *telemetry.Telemetry, logger *logrus.Logger) *Publisher {
	if logger == nil {
		logger = logrus.New()
	}

	fmt.Println("BROKERS:", cfg.Brokers)
	fmt.Println("USERNAME:", cfg.Username)
	fmt.Println("PASSWORD LENGTH:", len(cfg.Password))

	mechanism := kafkaconfig.NewMechanism(cfg)
	transport, err := kafkaconfig.NewTransport(cfg)
	if err != nil {
		panic(err)
	}

	return &Publisher{
		writer: &kafka.Writer{
			Addr:      kafka.TCP(cfg.Brokers...),
			Transport: transport,
			Balancer:  &kafka.LeastBytes{},
			Async:     false,
		},
		mechanism: mechanism,
		tlsConfig: transport.TLS,
		brokers:   cfg.Brokers,
		telemetry: telemetry,
		logger:    logger,
	}
}

func (p *Publisher) Publish(ctx context.Context, topic string, key string, event interface{}) error {
	ctx, span := p.telemetry.Start(ctx, "Publisher.Publish")
	defer span.End()

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	var headers []kafka.Header

	for k, v := range carrier {
		headers = append(headers,
			kafka.Header{
				Key:   k,
				Value: []byte(v),
			})
	}

	value, err := Encode(event)
	if err != nil {
		p.telemetry.RecordError(span, err)
		return err
	}

	if len(value) == 0 {
		p.telemetry.RecordError(span, fmt.Errorf("empty kafka message for topic=%s key=%s", topic, key))
		return fmt.Errorf("empty kafka message for topic=%s key=%s", topic, key)
	}

	p.writer.Logger = kafka.LoggerFunc(func(msg string, args ...interface{}) {
		fmt.Printf("[KAFKA] "+msg+"\n", args...)
	})

	ctx, kafkaSpan := p.telemetry.Start(ctx, "Kafka.WriteMessages")
	defer kafkaSpan.End()

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic:   topic,
		Key:     []byte(key),
		Value:   value,
		Time:    time.Now(),
		Headers: headers,
	})

	return err
}

func (p *Publisher) Ping(ctx context.Context) error {
	dialer := &kafka.Dialer{
		SASLMechanism: p.mechanism,
		TLS:           p.tlsConfig,
	}

	conn, err := dialer.DialContext(ctx, "tcp", p.brokers[0])
	if err != nil {
		return err
	}

	defer conn.Close()

	_, err = conn.Brokers()
	return err
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
