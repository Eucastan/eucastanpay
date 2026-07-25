package routes

import (
	commonmw "github.com/Eucastan/eucastanpay/common/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuth(api *gin.RouterGroup, deps Dependencies, h Handlers) {

	public := api.Group("")

	{
		auth := public.Group("/auth")

		auth.POST("/register", h.User.Register)
		auth.POST("/login", h.User.Login)
	}

	admin := api.Group("/admin/auth")
	{
		admin.POST("/register", h.Admin.CreateBootstrapAdmin)
		admin.POST("/login", h.Admin.Login)
	}

	_ = commonmw.Auth
}
