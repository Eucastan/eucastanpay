package bootstrap

import "github.com/Eucastan/eucastanpay/services/api-gateway/internal/application/service"

type Applications struct {
	userApp     *service.UserApplication
	adminApp    *service.AdminApplication
	accountApp  *service.AccountApplication
	transferApp *service.TransferApplication
	ledgerApp   *service.LedgerApplication
	auditApp    *service.AuditApplication
}

func (a *App) initApplications() {

	a.applications = &Applications{
		userApp:     service.NewUserApplication(a.gateways.userGateway),
		adminApp:    service.NewAdminApplication(a.gateways.adminGateway),
		accountApp:  service.NewAccountApplication(a.gateways.accountGateway),
		transferApp: service.NewTransferApplication(a.gateways.transferGateway),
		ledgerApp:   service.NewLedgerApplication(a.gateways.ledgerGateway),
		auditApp:    service.NewAuditApplication(a.gateways.auditGateway),
	}
}
