package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

const RequestContextKey = "request_context"

func RequestContext(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(
			c.Request.Context(),
			timeout,
		)

		c.Set(RequestContextKey, ctx)

		defer cancel()

		c.Next()
	}
}
