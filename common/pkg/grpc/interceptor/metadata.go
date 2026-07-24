package interceptor

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	UserIDKey        = "x-user-id"
	RoleKey          = "x-user-role"
	AdminIDKey       = "x-admin-id"
	AdminRoleKey     = "x-admin-role"
	CorrelationKey   = "x-correlation-id"
	AuthorizationKey = "authorization"
	RequestIDKey     = "x-request-id"
)

func MetadataFromContext(ctx context.Context) metadata.MD {

	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		return md
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return md
	}

	return metadata.New(nil)
}

func OutgoingContext(ctx context.Context, md metadata.MD) context.Context {
	return metadata.NewOutgoingContext(
		ctx,
		md,
	)
}
