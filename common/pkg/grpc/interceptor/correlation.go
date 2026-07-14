package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func CorrelationClientInterceptor() grpc.UnaryClientInterceptor {

	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		incoming := MetadataFromContext(ctx)
		outgoing := metadata.New(nil)

		copyKey := func(key string) {
			values := incoming.Get(key)
			if len(values) > 0 {
				outgoing.Set(key, values...)
			}
		}

		copyKey(RequestIDKey)
		copyKey(CorrelationKey)
		copyKey(UserIDKey)
		copyKey(RoleKey)

		ctx = metadata.NewOutgoingContext(ctx, outgoing)

		return invoker(ctx, method, req, reply, cc, opts...)

	}

}
