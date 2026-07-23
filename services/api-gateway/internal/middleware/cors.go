package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {

	return func(c *gin.Context) {

		h := c.Writer.Header()

		h.Set("Access-Control-Allow-Origin", "*")

		h.Set("Access-Control-Allow-Headers",
			"Authorization, Content-Type, X-Request-ID, X-Correlation-ID",
		)

		h.Set("Access-Control-Allow-Methods",
			"GET,POST,PUT,PATCH,DELETE,OPTIONS",
		)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
