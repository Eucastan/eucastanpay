package interceptor

import (
	"context"
	"strings"

	"github.com/Eucastan/eucastanpay/common/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(cfg string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata in context")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		const prefix = "Bearer "

		if !strings.HasPrefix(authHeader[0], prefix) {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header")
		}

		token := strings.TrimPrefix(authHeader[0], prefix)

		claims, err := auth.ValidateToken(token, cfg)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Propagate useful values into context
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "role", claims.Role)
		ctx = context.WithValue(ctx, "jwt_token", token)

		// Propagation for outgoing calls
		newMetadata := metadata.New(map[string]string{
			"user_id":       claims.UserID,
			"email":         claims.Email,
			"role":          claims.Role,
			"authorization": "Bearer " + token,
		})

		outgoingMD, _ := metadata.FromOutgoingContext(ctx)
		newMetadata = metadata.Join(outgoingMD, newMetadata)

		ctx = metadata.NewOutgoingContext(ctx, newMetadata)
		return handler(ctx, req)
	}
}
