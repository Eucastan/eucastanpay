package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/segmentio/kafka-go/sasl/plain"

	"github.com/Eucastan/eucastanpay/common/pkg/config"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/kafkaconfig"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type HandlerFunc func(ctx context.Context, msg []byte) error

type Consumer struct {
	brokers   []string
	username  string
	password  string
	groupID   string
	handlers  map[string]HandlerFunc
	readers   []*kafka.Reader
	mechanism *plain.Mechanism
	dialer    *kafka.Dialer
	mu        sync.Mutex
	wg        sync.WaitGroup
	telemetry *telemetry.Telemetry
	logger    *logrus.Logger
}

func NewConsumer(cfg config.KafkaConfig, groupID string, telemetry *telemetry.Telemetry, logger *logrus.Logger) *Consumer {
	if logger == nil {
		logger = logrus.New()
	}

	logger.Infof("BROKERS=%v", cfg.Brokers)
	logger.Infof("USERNAME=%s", cfg.Username)
	logger.Infof("TLS=%v", cfg.TLS)
	logger.Infof("SASL=%v", cfg.SASL)

	mechanism := kafkaconfig.NewMechanism(cfg)
	dialer, err := kafkaconfig.NewDialer(cfg)
	if err != nil {
		panic(err)
	}

	return &Consumer{
		brokers:   cfg.Brokers,
		username:  cfg.Username,
		password:  cfg.Password,
		groupID:   groupID,
		handlers:  make(map[string]HandlerFunc),
		readers:   []*kafka.Reader{},
		mechanism: mechanism,
		dialer:    dialer,
		telemetry: telemetry,
		logger:    logger,
	}
}

func (c *Consumer) Register(topic string, handler HandlerFunc) {
	c.logger.Infof("REGISTERING %s", topic)
	c.handlers[topic] = handler
}

// Start listening to multiple topics
func (c *Consumer) Start(ctx context.Context) {
	c.logger.Infof("registered handlers=%d", len(c.handlers))

	for topic, handler := range c.handlers {
		c.logger.Infof("starting consumer for %s", topic)
		c.wg.Add(1)

		go func(topic string, handler HandlerFunc) {
			defer c.wg.Done()
			c.consumeTopic(ctx, topic, handler)
		}(topic, handler)
	}
}

func (c *Consumer) consumeTopic(
	ctx context.Context,
	topic string,
	handler HandlerFunc,
) {
	ctx, span := c.telemetry.Start(ctx, "Consumer.consumeTopic")
	defer span.End()

	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:           c.brokers,
			GroupID:           c.groupID,
			Topic:             topic,
			Dialer:            c.dialer,
			MinBytes:          1,
			MaxBytes:          10e6,
			MaxWait:           time.Second,
			CommitInterval:    0,
			SessionTimeout:    30 * time.Second,
			HeartbeatInterval: 3 * time.Second,
			//StartOffset:       kafka.FirstOffset,
		},
	)
	c.mu.Lock()
	c.readers = append(c.readers, reader)
	c.mu.Unlock()

	c.logger.Infof("%+v", reader.Config())
	stats := reader.Stats()
	c.logger.Infof("%+v", stats)

	c.logger.Infof("consumer started topic=%s", topic)
	c.logger.Infof("Consumer topic=%s brokers=%v group=%s", topic, c.brokers, c.groupID)

	for {

		select {
		case <-ctx.Done():
			return
		default:
		}

		msgCtx, span := c.telemetry.Start(ctx, "FetchMessage")

		c.logger.Info("Consumer alive")
		c.logger.Infof("Consumer topic=%s brokers=%v group=%s", topic, c.brokers, c.groupID)
		c.logger.Infof("Calling FetchMessage for %s", topic)
		msg, err := reader.FetchMessage(msgCtx)
		if err != nil {
			span.End()
			if ctx.Err() != nil {
				return
			}

			c.logger.Printf("fetch error topic=%s err=%v\n", topic, err)
			continue
		}

		c.logger.Infof("FetchMessage returned for %s", topic)

		carrier := propagation.MapCarrier{}
		for _, h := range msg.Headers {
			carrier[h.Key] = string(h.Value)
		}

		msgCtx = otel.GetTextMapPropagator().Extract(msgCtx, carrier)
		c.logger.Infof("HANDLER START topic=%s key=%s", topic, string(msg.Key))

		if err := handler(msgCtx, msg.Value); err != nil {
			c.telemetry.RecordError(span, err)
			span.End()
			c.logger.Errorf("HANDLER FAILED topic=%s err=%v payload=%s", topic, err, string(msg.Value))
			continue
		}

		if err := reader.CommitMessages(ctx, msg); err != nil {
			if ctx.Err() != nil {
				span.End()
				return
			}

			c.telemetry.RecordError(span, err)
			c.logger.Printf("commit error topic=%s err=%v\n", topic, err)
		}
		span.End()
	}
}

func (c *Consumer) Close() error {

	c.mu.Lock()
	readers := append([]*kafka.Reader(nil), c.readers...)
	c.mu.Unlock()

	var firstErr error

	for _, r := range readers {
		if err := r.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	c.wg.Wait()

	return firstErr
}
