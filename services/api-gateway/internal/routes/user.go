package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(api *gin.RouterGroup, deps Dependencies, h Handlers) {

	users := UserGroup(api, deps.Config)

	{
		users.GET("/users", h.User.GetAllUsers)
		users.GET("/user", h.User.GetUser)
		users.PUT("/users/:id", h.User.UpdateUser)
		users.DELETE("/users/:id", h.User.DeleteUser)

		users.POST("/users/:user_id/kyc", h.User.CreateKYC)
		users.GET("/users/:user_id/kyc", h.User.GetKYC)
	}
}
