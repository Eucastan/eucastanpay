package interceptor

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	UserIDKey = "x-user-id"

	RoleKey = "x-user-role"

	CorrelationKey = "x-correlation-id"

	RequestIDKey = "x-request-id"
)

func MetadataFromContext(ctx context.Context) metadata.MD {

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return metadata.New(nil)
	}

	return md
}

func OutgoingContext(ctx context.Context, md metadata.MD) context.Context {
	return metadata.NewOutgoingContext(
		ctx,
		md,
	)
}
