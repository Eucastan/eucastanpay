package api

import (
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/ledger/config"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, h *handler.LedgerHandler, cfg *config.Config) {
	auth := r.Group("/api/v1")
	auth.Use(middleware.Auth(cfg.SharedCfg.JWTSecret))
	{
		auth.GET("/ledgers", middleware.RequireRole("user", "super_admin"), h.GetAllLedgers)
		auth.GET("/ledgers/:id", middleware.RequireRole("user", "super_admin"), h.GetLedger)
		auth.GET("/accounts/:account_id/balance", middleware.RequireRole("user", "super_admin"), h.GetAccountBalance)
		auth.GET("/accounts/:account_id/reconcile", middleware.RequireRole("user", "super_admin"), h.ReconciliationResult)
	}
}
