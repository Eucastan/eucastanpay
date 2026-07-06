package api

import (
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/ledger/config"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, h *handler.LedgerHandler, cfg *config.Config) {
	v1 := r.Group("/api/v1")
	v1.Use(middleware.Auth(cfg.JWTSecret))
	{
		v1.GET("/ledgers/:id", h.GetLedger)
		v1.GET("/ledgers/entry_type", h.GetLedgerByEntry)
		v1.GET("/accounts/:account_id/balance", h.GetAccountBalance)
	}
}
