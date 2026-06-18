package consumer

import (
	"context"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
)

func RetryHandler(
	handler HandlerFunc,
	publisher *producer.Publisher,
	topic string,
	dqlTopic string,
	maxRetries int,
) HandlerFunc {

	return func(ctx context.Context, msg []byte) error {
		var err error

		var count int
		for i := 0; i < maxRetries; i++ {
			err = handler(ctx, msg)
			if err == nil {
				return nil
			}

			count = i
			time.Sleep(time.Duration(i+1) * time.Second) // backoff
		}

		// send to DLQ after retries exhausted
		_ = publisher.Publish(ctx, dqlTopic, topic, events.DLQEvent{
			OriginalTopic: topic,
			Error:         err.Error(),
			RetryCount:    count,
			Payload:       string(msg),
			FailedAt:      time.Now().Unix(),
		})

		return err
	}
}
