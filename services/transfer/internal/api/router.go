package api

import (
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/transfer/config"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, h *handler.TransferHandler, cfg *config.Config) {
	v1 := r.Group("/api/v1")
	v1.Use(middleware.Auth(cfg.JWTSecret))
	{
		v1.POST("/transfers", middleware.RequireRole("user", "super_admin"), h.TransferFromUser)
		v1.GET("/transfers/:id", h.GetTransfer)
	}
}
