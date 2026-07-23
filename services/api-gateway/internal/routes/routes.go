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

	RegisterHealth(r, healthcheck.New(
		deps.Config.ServiceName,
		deps.Config.Version,
		deps.Logger,
	),
	)

	RegisterSwaggerRoutes(r)
	RegisterAuth(r, deps, h)

	RegisterUserRoutes(r, deps, h)
	RegisterAdminRoutes(r, deps, h)
	RegisterAccountRoutes(r, deps, h)
	RegisterTransferRoutes(r, deps, h)
	RegisterLedgerRoutes(r, deps, h)
	RegisterAuditRoutes(r, deps, h)
}
