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

		for i := 0; i <= maxRetry; i++ {

			err = invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}

			if i == maxRetry {
				break
			}

			time.Sleep(backoff)
			backoff *= 2
		}

		return err
	}

}
