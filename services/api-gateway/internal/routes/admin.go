package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(api *gin.RouterGroup, deps Dependencies, h Handlers) {

	admin := AdminGroup(api, deps.Config)

	{
		admin.GET("/", h.Admin.ListAdmins)
		admin.GET("/:id", h.Admin.GetAdmin)

		admin.POST("/", h.Admin.CreateAdmin)
		admin.PUT("/:id", h.Admin.UpdateAdmin)
		admin.DELETE("/:id", h.Admin.DeleteAdmin)
		admin.PATCH("/:user_id/kyc", h.User.ApproveKYC)
	}
}
