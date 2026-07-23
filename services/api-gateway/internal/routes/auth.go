package routes

import (
	commonmw "github.com/Eucastan/eucastanpay/common/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuth(r *gin.Engine, deps Dependencies, h Handlers) {

	public := r.Group("/api/v1")

	{
		auth := public.Group("/auth")

		auth.POST("/register", h.User.Register)
		auth.POST("/login", h.User.Login)

		admin := public.Group("/admin/auth")

		admin.POST("/register", h.Admin.CreateBootstrapAdmin)
		admin.POST("/login", h.Admin.Login)
	}

	_ = commonmw.Auth
}
