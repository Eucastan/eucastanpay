package routes

import (
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterTransferRoutes(api *gin.RouterGroup, deps Dependencies, h Handlers) {

	admin := AdminGroup(api, deps.Config)

	{
		admin.GET("/transfers", h.Transfer.GetAllTransfers)
		admin.GET("/transfers/:id", h.Transfer.GetTransfer)
		admin.GET("/accounts/:account_id", h.Transfer.GetTransferByAccID)
		admin.GET("/accounts/:account_id/transfer", h.Transfer.ReconcileAccount)
		admin.POST("/transfers/:reference", middleware.RequireIdempotencyKey(), h.Transfer.ReverseTransfer)
	}

	user := UserGroup(api, deps.Config)

	{
		user.POST("/transfers", middleware.RequireIdempotencyKey(), h.Transfer.TransferFromUser)
		user.GET("/transfers/me", h.Transfer.GetTransferByUserID)
	}
}
