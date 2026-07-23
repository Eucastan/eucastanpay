package routes

import (
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterTransferRoutes(r *gin.Engine, deps Dependencies, h Handlers) {

	admin := AdminGroup(r, deps.Config)

	{
		admin.GET("/transfers", h.Transfer.GetAllTransfers)
		admin.GET("/transfers/:id", h.Transfer.GetTransfer)
		admin.POST("/accounts/:account_id/reconcile", h.Transfer.ReconcileAccount)
		admin.POST("/transfers/:reference", middleware.RequireIdempotencyKey(), h.Transfer.ReverseTransfer)
	}

	user := UserGroup(r, deps.Config)

	{
		user.POST("/transfers", middleware.RequireIdempotencyKey(), h.Transfer.TransferFromUser)
		user.GET("/transfers/me", h.Transfer.GetTransferByUserID)
	}
}
