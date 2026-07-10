package api

import (
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/notification/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, h *handler.NotificationHandler) {
	auth := r.Group("/api/v1")
	{
		auth.GET("/notifications", middleware.RequireRole("super_admin", "admin", "user"), h.GetNotifications)
	}
}
