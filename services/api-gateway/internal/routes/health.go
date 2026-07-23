package routes

import (
	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"

	"github.com/gin-gonic/gin"
)

func RegisterHealth(router *gin.Engine, health *healthcheck.Checker) {

	router.GET("/health", health.Health)
	router.GET("/ready", health.Readiness)
	router.GET("/live", health.Liveness)
}
