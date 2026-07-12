package api

import (
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/transfer/config"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, h *handler.TransferHandler, cfg *config.Config) {
	auth := r.Group("/api/v1")
	auth.Use(middleware.Auth(cfg.SharedCfg.JWTSecret))
	{
		auth.POST("/transfers", middleware.RequireRole("user", "super_admin"), h.TransferFromUser)
		auth.POST("/transfers/:reference", middleware.RequireRole("user", "super_admin"), h.ReverseTransfer)
		auth.GET("/transfers", middleware.RequireRole("user", "super_admin"), h.GetAllTransfer)
		auth.GET("/transfers/:id", middleware.RequireRole("user", "super_admin"), h.GetTransfer)
		auth.POST("/account/:account_id/reconcile", middleware.RequireRole("user", "super_admin"), h.ReconcileAccount)
	}
}
