package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(brokers []string) *Publisher {
	return &Publisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
			Async:    false, // Set true in high-throughput scenarios
		},
	}
}

func (p *Publisher) Publish(ctx context.Context, topic string, key string, event interface{}) error {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if len(value) == 0 {
		return fmt.Errorf("empty kafka message for topic=%s key=%s", topic, key)
	}

	p.writer.Logger = kafka.LoggerFunc(func(msg string, args ...interface{}) {
		fmt.Printf("[KAFKA] "+msg+"\n", args...)
	})

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	})
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
