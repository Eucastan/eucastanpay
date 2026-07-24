package proxy

import (
	"context"

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

	keys := []string{
		"request_id",
		"correlation_id",
		"user_id",
		"admin_id",
		"email",
		"role",
		"admin_role",
	}

	for _, key := range keys {
		if value := c.GetString(key); value != "" {
			md.Set(key, value)
		}
	}

	if token := c.GetString("token"); token != "" {
		md.Set("authorization", "Bearer "+token)
	}

	return metadata.NewOutgoingContext(baseCtx, md)
}
