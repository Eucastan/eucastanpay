package proxy

import (
	"context"

	commongrpc "github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

func Context(c *gin.Context) context.Context {

	baseCtx := c.Request.Context()

	if ctx, ok := c.Get(middleware.RequestContextKey); ok {
		baseCtx = ctx.(context.Context)
	}

	md := metadata.New(nil)

	if v := c.GetString("request_id"); v != "" {
		md.Set(commongrpc.RequestIDKey, v)
	}

	if v := c.GetString("correlation_id"); v != "" {
		md.Set(commongrpc.CorrelationKey, v)
	}

	if v := c.GetString("user_id"); v != "" {
		md.Set(commongrpc.UserIDKey, v)
	}

	if v := c.GetString("role"); v != "" {
		md.Set(commongrpc.RoleKey, v)
	}

	if v := c.GetString("admin_id"); v != "" {
		md.Set(commongrpc.AdminIDKey, v)
	}

	if v := c.GetString("admin_role"); v != "" {
		md.Set(commongrpc.AdminRoleKey, v)
	}

	if token := c.GetString("token"); token != "" {
		md.Set(commongrpc.AuthorizationKey, "Bearer "+token)
	}

	return commongrpc.OutgoingContext(baseCtx, md)
}
