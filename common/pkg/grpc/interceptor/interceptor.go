package interceptor

import (
	"context"
	"log"
	"strings"

	"github.com/Eucastan/eucastanpay/common/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var publicMethods = map[string]bool{
	"/user.UserService/Register":        true,
	"/user.UserService/Login":           true,
	"/user.UserService/VerifyEmail":     true,
	"/user.UserService/RefreshToken":    true,
	"/user.UserService/RequestPassword": true,
}

func AuthInterceptor(cfg string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("========================")
		log.Printf("METHOD: %s", info.FullMethod)

		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata in context")
		}

		log.Printf("Incoming metadata: %+v\n", md)
		log.Printf("HAS METADATA: %v", ok)
		log.Printf("METADATA: %#v", md)

		authHeader := md.Get("authorization")
		log.Printf("AUTH HEADER: %#v", authHeader)
		log.Println("========================")
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
			UserIDKey:        claims.UserID,
			RoleKey:          claims.Role,
			AuthorizationKey: "Bearer " + token,
		})

		outgoingMD, _ := metadata.FromOutgoingContext(ctx)
		newMetadata = metadata.Join(outgoingMD, newMetadata)

		ctx = metadata.NewOutgoingContext(ctx, newMetadata)
		return handler(ctx, req)
	}
}
