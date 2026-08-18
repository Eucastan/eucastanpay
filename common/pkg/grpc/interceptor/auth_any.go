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

	// "/admin.AdminService/Register":     true,
	"/admin.AdminService/BootstrapAdmin": true,
	"/admin.AdminService/Login":          true,
	"/admin.AdminService/RefreshToken":   true,
}

func AuthAnyInterceptor(secret string) grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get(AuthorizationKey)
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		const prefix = "Bearer "

		if !strings.HasPrefix(authHeader[0], prefix) {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header")
		}

		token := strings.TrimPrefix(authHeader[0], prefix)

		// Try USER token
		if claims, err := auth.ValidateToken(token, secret); err == nil {

			log.Printf("Authenticated USER %s", claims.UserID)

			ctx = context.WithValue(ctx, ContextPrincipal, PrincipalUser)
			ctx = context.WithValue(ctx, ContextUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextRole, claims.Role)
			ctx = context.WithValue(ctx, ContextJWTToken, token)

			outgoing := metadata.New(map[string]string{
				UserIDKey:        claims.UserID,
				RoleKey:          claims.Role,
				AuthorizationKey: prefix + token,
			})

			if md, ok := metadata.FromOutgoingContext(ctx); ok {
				outgoing = metadata.Join(md, outgoing)
			}

			ctx = metadata.NewOutgoingContext(ctx, outgoing)

			return handler(ctx, req)
		}

		// Try ADMIN token
		if claims, err := auth.ValidateAdminToken(token, secret); err == nil {

			log.Printf("Authenticated ADMIN %s", claims.AdminID)

			ctx = context.WithValue(ctx, ContextPrincipal, PrincipalAdmin)
			ctx = context.WithValue(ctx, ContextAdminID, claims.AdminID)
			ctx = context.WithValue(ctx, ContextAdminRole, claims.Role)
			ctx = context.WithValue(ctx, ContextJWTToken, token)

			outgoing := metadata.New(map[string]string{
				AdminIDKey:       claims.AdminID,
				AdminRoleKey:     claims.Role,
				AuthorizationKey: prefix + token,
			})

			if md, ok := metadata.FromOutgoingContext(ctx); ok {
				outgoing = metadata.Join(md, outgoing)
			}

			ctx = metadata.NewOutgoingContext(ctx, outgoing)

			return handler(ctx, req)
		}

		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
}
