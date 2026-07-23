package middleware

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func Logger(log *logrus.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()
		c.Next()

		log.WithFields(logrus.Fields{

			"method":         c.Request.Method,
			"path":           c.Request.URL.Path,
			"status":         c.Writer.Status(),
			"latency":        time.Since(start),
			"request_id":     c.GetString("request_id"),
			"correlation_id": c.GetString("correlation_id"),
		}).Info("gateway request")
	}
}
