package proxy

import (
	"context"

	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Context(c *gin.Context) context.Context {

	ctx, ok := c.Get(
		middleware.OutgoingContextKey,
	)

	if !ok {
		return c.Request.Context()
	}

	return ctx.(context.Context)
}
