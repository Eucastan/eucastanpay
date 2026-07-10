package api

import (
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/audit/config"
	"github.com/Eucastan/eucastanpay/services/audit/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, h *handler.AuditHandler, cfg *config.Config) {
	auth := r.Group("/api/v1")
	auth.Use(middleware.Auth(cfg.JWTSecret))
	{
		auth.GET("/audit/search", middleware.RequireRole("user", "admin", "super_admin"), h.SearchAuditLogs)
		auth.GET("/audit/:id", middleware.RequireRole("user", "admin", "super_admin"), h.GetAuditRead)
	}
}
