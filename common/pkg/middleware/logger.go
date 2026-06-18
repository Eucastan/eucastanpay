package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger(log *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		log.WithFields(logrus.Fields{
			"method":  ctx.Request.Method,
			"path":    ctx.Request.URL.Path,
			"status":  ctx.Writer.Status(),
			"start":   start,
			"latency": time.Since(start),
			"ip":      ctx.ClientIP(),
			"user_id": ctx.GetString("user_id"),
		}).Info("Request processed")
	}
}
