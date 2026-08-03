package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func RetryInterceptor(maxRetry int) grpc.UnaryClientInterceptor {

	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		var err error

		backoff := 200 * time.Millisecond

		for attempt := 0; attempt <= maxRetry; attempt++ {

			err = invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}

			if !shouldRetry(err) {
				return err
			}

			if attempt == maxRetry {
				return err
			}

			select {
			case <-time.After(backoff):

			case <-ctx.Done():
				return ctx.Err()
			}

			backoff *= 2
		}

		return err
	}
}
