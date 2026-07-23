package bootstrap

import "github.com/gin-gonic/gin"

func (a *App) initRouter() {
	a.router = gin.Default()

	a.router.GET("/health", a.health.Health)
	a.router.GET("/live", a.health.Liveness)
	a.router.GET("/ready", a.health.Readiness)
}
