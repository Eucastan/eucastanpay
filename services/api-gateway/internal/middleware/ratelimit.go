package middleware

import (
	"net/http"

	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/ratelimiter"
	"github.com/gin-gonic/gin"
)

func RateLimit(
	limiter ratelimiter.Limiter,
	limit int,
	windowSeconds int,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx := c.Request.Context()

		key := buildKey(c)

		ok, err := limiter.Allow(
			ctx,
			key,
			limit,
			windowSeconds,
		)

		if err != nil {

			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"success": false,
					"message": "rate limiter unavailable",
				},
			)

			return
		}

		if !ok {

			c.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				gin.H{
					"success": false,
					"message": "too many requests",
				},
			)

			return
		}

		c.Next()
	}
}

func buildKey(c *gin.Context) string {

	if admin := c.GetString("admin_id"); admin != "" {
		return "admin:" + admin
	}

	if user := c.GetString("user_id"); user != "" {
		return "user:" + user
	}

	return "ip:" + c.ClientIP()
}
