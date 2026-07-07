package interceptor

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func AppendJWTToContext(ctx context.Context, tokenStr string) context.Context {
	return metadata.AppendToOutgoingContext(
		ctx,
		"authorization",
		"Bearer "+tokenStr,
	)
}
