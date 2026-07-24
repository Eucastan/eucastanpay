package routes

import (
	"time"

	sharedmw "github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RegisterGlobalMiddleware(r *gin.Engine, logger *logrus.Logger) {

	r.Use(
		middleware.Recovery(),
		middleware.RequestID(),
		sharedmw.CorrelationMiddleware(),
		middleware.CORS(),
		middleware.RequestContext(10*time.Second),
		middleware.Logger(logger),
	)
}
