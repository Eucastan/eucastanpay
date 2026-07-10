package api

import (
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/account/config"
	"github.com/Eucastan/eucastanpay/services/account/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, h *handler.AccountHandler, cfg *config.Config) {

	auth := r.Group("/api/v1")
	{
		auth.POST("/account", middleware.RequireRole("user", "super_admin", "admin"), h.OpenAccount)
		auth.GET("/account/:id", h.GetUserAccount)
		auth.GET("/account/:id/balance", h.GetBalance)
		auth.PUT("/account/:id/pay-in", h.InitiatePayIn)
		auth.PUT("/account/:id/withdraw", h.WithDrawal)
	}
}
