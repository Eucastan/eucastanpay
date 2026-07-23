package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/clients"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/gateway"
)

type Gateways struct {
	userGateway     *gateway.UserGateway
	adminGateway    *gateway.AdminGateway
	accountGateway  *gateway.AccountGateway
	transferGateway *gateway.TransferGateway
	ledgerGateway   *gateway.LedgerGateway
	auditGateway    *gateway.AuditGateway
}

func (a *App) initGateways() {

	userClient := clients.User(a.manager)
	adminClient := clients.Admin(a.manager)
	accountClient := clients.Account(a.manager)
	transferClient := clients.Transfer(a.manager)
	ledgerClient := clients.Ledger(a.manager)
	auditClient := clients.Audit(a.manager)

	a.gateways = &Gateways{
		userGateway:     gateway.NewUserGateway(userClient),
		adminGateway:    gateway.NewAdminGateway(adminClient),
		accountGateway:  gateway.NewAccountGateway(accountClient),
		transferGateway: gateway.NewTransferGateway(transferClient),
		ledgerGateway:   gateway.NewLedgerGateway(ledgerClient),
		auditGateway:    gateway.NewAuditGateway(auditClient),
	}
}
