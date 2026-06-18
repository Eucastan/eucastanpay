package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Recovery(log *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(logrus.Fields{
					"error": err,
					"stack": string(debug.Stack()),
				}).Error("panic recovered")

				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
			}
		}()

		ctx.Next()
	}
}
