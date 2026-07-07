package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
)

func RetryHandler(
	handler HandlerFunc,
	publisher *producer.Publisher,
	topic string,
	dlqTopic string,
	t *telemetry.Telemetry,
	maxRetries int,
) HandlerFunc {

	return func(ctx context.Context, msg []byte) error {
		ctx, span := t.Start(ctx, "RetryHandler")
		defer span.End()

		var handlerErr error

		for i := 0; i < maxRetries; i++ {
			handlerErr = handler(ctx, msg)
			if handlerErr == nil {
				return nil
			}

			time.Sleep(time.Duration(i+1) * time.Second) // backoff
		}

		// send to DLQ after retries exhausted
		dlqErr := publisher.Publish(ctx, dlqTopic, topic, events.DLQEvent{
			OriginalTopic: topic,
			Error:         handlerErr.Error(),
			RetryCount:    maxRetries,
			Payload:       string(msg),
			FailedAt:      time.Now().Unix(),
		})
		if dlqErr != nil {
			return fmt.Errorf(
				"handler failed: %v, dlq failed: %w",
				handlerErr,
				dlqErr,
			)
		}

		return handlerErr
	}
}
