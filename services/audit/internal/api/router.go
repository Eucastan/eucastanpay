package api

import (
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/audit/config"
	"github.com/Eucastan/eucastanpay/services/audit/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, h *handler.AuditHandler, cfg *config.Config) {
	v1 := r.Group("/api/v1")
	v1.Use(middleware.Auth(cfg.JWTSecret))
	{
		v1.GET("/audit/search", h.Search)
	}
}
