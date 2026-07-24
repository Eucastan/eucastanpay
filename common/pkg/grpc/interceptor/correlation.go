package interceptor

import (
	"context"
	"log"

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
		log.Printf("CLIENT INTERCEPTOR CALLED: %s", method)

		incoming := MetadataFromContext(ctx)
		outgoing := metadata.New(nil)

		log.Printf("OUTGOING BEFORE COPY: %+v", incoming)

		copyKey := func(key string) {
			values := incoming.Get(key)
			if len(values) > 0 {
				outgoing.Set(key, values...)
			}
		}

		copyKey(RequestIDKey)
		copyKey(CorrelationKey)
		copyKey(UserIDKey)
		copyKey(AuthorizationKey)
		copyKey(RoleKey)
		copyKey(AdminIDKey)
		copyKey(AdminRoleKey)

		ctx = metadata.NewOutgoingContext(ctx, outgoing)
		// Debug below
		md, _ := metadata.FromOutgoingContext(ctx)
		log.Printf("OUTGOING AFTER COPY: %+v", md)

		return invoker(ctx, method, req, reply, cc, opts...)

	}

}
