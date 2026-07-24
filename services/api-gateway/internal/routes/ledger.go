package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterLedgerRoutes(api *gin.RouterGroup, deps Dependencies, h Handlers) {

	admin := AdminGroup(api, deps.Config)

	{
		admin.GET("/ledgers", h.Ledger.GetAllLedgers)
		admin.GET("/ledgers/:id", h.Ledger.GetLedger)
		admin.GET("/ledgers/account/:account_id", h.Ledger.GetLedgerByAccountID)
		admin.GET("/accounts/:account_id/balance", h.Ledger.GetAccountBalance)
		admin.GET("/accounts/:account_id/reconcile", h.Ledger.ReconciliationResult)
	}

	user := UserGroup(api, deps.Config)

	{
		user.GET("/ledgers/me", h.Ledger.GetLedgerByUserID)
	}
}
