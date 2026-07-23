package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

const OutgoingContextKey = "grpc_outgoing_ctx"

func OutgoingContext(timeout time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(
			c.Request.Context(),
			timeout,
		)

		md := metadata.New(nil)

		if v := c.GetString("request_id"); v != "" {
			md.Set("request-id", v)
		}

		if v := c.GetString("correlation_id"); v != "" {
			md.Set("correlation-id", v)
		}

		if v := c.GetString("user_id"); v != "" {
			md.Set("user-id", v)
		}

		if v := c.GetString("email"); v != "" {
			md.Set("email", v)
		}

		if v := c.GetString("role"); v != "" {
			md.Set("role", v)
		}

		outgoing := metadata.NewOutgoingContext(
			ctx,
			md,
		)

		c.Set(
			OutgoingContextKey,
			outgoing,
		)

		defer cancel()

		c.Next()
	}
}
