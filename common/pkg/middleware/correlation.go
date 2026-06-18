package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const CorrelationIDKey = "correlation_id"

func CorrelationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get existing correlation ID from header (for tracing across services)
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		// Add to Gin context
		c.Set(CorrelationIDKey, correlationID)

		// Add to Go context for downstream use
		ctx := context.WithValue(c.Request.Context(), CorrelationIDKey, correlationID)
		c.Request = c.Request.WithContext(ctx)

		// Add to response header for client tracing
		c.Header("X-Correlation-ID", correlationID)

		c.Next()
	}
}
