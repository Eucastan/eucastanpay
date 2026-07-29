package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(api *gin.RouterGroup, deps Dependencies, h Handlers) {

	s := SuperAdmin(api, deps.Config)

	{
		s.GET("/", h.Admin.ListAdmins)
		s.GET("/:id", h.Admin.GetAdmin)

		s.POST("/", h.Admin.CreateAdmin)
		s.PATCH("/:id", h.Admin.UpdateAdmin)
		s.DELETE("/:id", h.Admin.DeleteAdmin)
	}

	admin := AdminGroup(api, deps.Config)
	{
		admin.GET("/me", h.Admin.GetAdminProfile)
	}
}
