package bootstrap

import (
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/routes"
	"github.com/gin-gonic/gin"
)

func (a *App) initRouter() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	a.router = r
}

func (a *App) registerRoutes() {

	routes.Register(
		a.router,
		routes.Dependencies{
			Config:      a.cfg,
			Logger:      a.logger,
			RateLimiter: a.rateLimiter,
		},

		routes.Handlers{
			User:     a.handlers.User,
			Admin:    a.handlers.Admin,
			Account:  a.handlers.Account,
			Transfer: a.handlers.Transfer,
			Ledger:   a.handlers.Ledger,
			Audit:    a.handlers.Audit,
		},
	)
}
