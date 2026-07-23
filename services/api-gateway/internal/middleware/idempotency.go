package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const IdempotencyKeyContext = "idempotency_key"

func RequireIdempotencyKey() gin.HandlerFunc {

	return func(c *gin.Context) {

		key := c.GetHeader("Idempotency-Key")

		if key == "" {

			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"success": false,
					"message": "Idempotency-Key header is required",
				},
			)

			return
		}

		c.Set(
			IdempotencyKeyContext,
			key,
		)

		c.Next()
	}
}
