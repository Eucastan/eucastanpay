package api

import (
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/account/config"
	"github.com/Eucastan/eucastanpay/services/account/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, h *handler.AccountHandler, cfg *config.Config) {
	r.Use(middleware.Auth(cfg.JWTSecret))
	v1 := r.Group("/api/v1")
	{
		v1.POST("/account", middleware.RequireRole("user", "super_admin", "admin"), h.OpenAccount)
		v1.GET("/account/:id/balance", h.GetBalance)
	}
}
