package routes

import (
	"time"

	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RegisterGlobalMiddleware(r *gin.Engine, logger *logrus.Logger) {

	r.Use(
		middleware.Recovery(),
		middleware.RequestID(),
		middleware.CORS(),
		middleware.RequestContext(10*time.Second),
		middleware.Logger(logger),
	)
}
