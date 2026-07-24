package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterAccountRoutes(api *gin.RouterGroup, deps Dependencies, h Handlers) {

	admin := AdminGroup(api, deps.Config)

	{
		admin.GET("/account", h.Account.GetAllUsersAccount)
		admin.GET("/account/:id/balance", h.Account.GetBalance)
		admin.PUT("/account/:id/pay-in", h.Account.InitiatePayIn)
		admin.PUT("/account/:id/withdraw", h.Account.WithDrawal)
	}

	user := UserGroup(api, deps.Config)

	{
		user.GET("/account/:id", h.Account.GetUserAccount)
	}
}
