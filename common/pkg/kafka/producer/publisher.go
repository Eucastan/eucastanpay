package producer

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go/sasl/plain"

	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Publisher struct {
	writer    *kafka.Writer
	telemetry *telemetry.Telemetry
}

func NewPublisher(brokers []string, username, password string, telemetry *telemetry.Telemetry) *Publisher {

	fmt.Println("BROKERS:", brokers)
	fmt.Println("USERNAME:", username)
	fmt.Println("PASSWORD LENGTH:", len(password))

	mechanism := plain.Mechanism{
		Username: username,
		Password: password,
	}

	transport := &kafka.Transport{
		SASL: mechanism,
		TLS:  &tls.Config{},
	}

	return &Publisher{
		writer: &kafka.Writer{
			Addr:      kafka.TCP(brokers...),
			Transport: transport,
			Balancer:  &kafka.LeastBytes{},
			Async:     false,
		},
		telemetry: telemetry,
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

func (p *Publisher) Close() error {
	return p.writer.Close()
}
