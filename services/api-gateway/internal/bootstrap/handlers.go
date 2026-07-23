package bootstrap

import (
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/handler"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/proxy"
)

type Handlers struct {
	User     *handler.UserHandler
	Admin    *handler.AdminHandler
	Account  *handler.AccountHandler
	Transfer *handler.TransferHandler
	Ledger   *handler.LedgerHandler
	Audit    *handler.AuditHandler
}

func (a *App) initHandlers() {

	base := proxy.NewBase(
		a.logger,
		a.telemetry,
	)

	a.handlers = &Handlers{
		User:     handler.NewUserHandler(base, a.applications.userApp),
		Admin:    handler.NewAdminHandler(base, a.applications.adminApp),
		Account:  handler.NewAccountHandler(base, a.applications.accountApp),
		Transfer: handler.NewTransferHandler(base, a.applications.transferApp),
		Ledger:   handler.NewLedgerHandler(base, a.applications.ledgerApp),
		Audit:    handler.NewAuditHandler(base, a.applications.auditApp),
	}
}
