package routes

import (
	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
	"github.com/Eucastan/eucastanpay/services/api-gateway/config"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/handler"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/ratelimiter"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	User     *handler.UserHandler
	Admin    *handler.AdminHandler
	Account  *handler.AccountHandler
	Transfer *handler.TransferHandler
	Ledger   *handler.LedgerHandler
	Audit    *handler.AuditHandler
}

type Dependencies struct {
	Logger      *logrus.Logger
	Config      *config.Config
	RateLimiter ratelimiter.Limiter
}

func Register(r *gin.Engine, deps Dependencies, h Handlers) {

	RegisterGlobalMiddleware(r, deps.Logger)

	RegisterHealth(r, healthcheck.New(
		deps.Config.ServiceName,
		deps.Config.Version,
		deps.Logger,
	),
	)

	api := r.Group("/api/v1")

	RegisterSwaggerRoutes(r)
	RegisterAuth(api, deps, h)

	RegisterUserRoutes(api, deps, h)
	RegisterAdminRoutes(api, deps, h)
	RegisterAccountRoutes(api, deps, h)
	RegisterTransferRoutes(api, deps, h)
	RegisterLedgerRoutes(api, deps, h)
	RegisterAuditRoutes(api, deps, h)
}
