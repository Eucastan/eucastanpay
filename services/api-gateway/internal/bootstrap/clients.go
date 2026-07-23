package bootstrap

import (
	"time"

	manager "github.com/Eucastan/eucastanpay/common/pkg/grpc"
)

func (a *App) initClients() error {

	user := manager.ServiceConfig{
		Name:     "user",
		Address:  a.cfg.UserGRPCAddr,
		Insecure: true,
		Timeout:  time.Duration(5 * time.Second),
		Retries:  3,
	}

	admin := manager.ServiceConfig{
		Name:     "admin",
		Address:  a.cfg.AdminGRPCAddr,
		Insecure: true,
		Timeout:  time.Duration(5 * time.Second),
		Retries:  3,
	}

	account := manager.ServiceConfig{
		Name:     "account",
		Address:  a.cfg.AccountGRPCAddr,
		Insecure: true,
		Timeout:  time.Duration(5 * time.Second),
		Retries:  3,
	}

	transfer := manager.ServiceConfig{
		Name:     "transfer",
		Address:  a.cfg.TransferGRPCAddr,
		Insecure: true,
		Timeout:  time.Duration(5 * time.Second),
		Retries:  3,
	}

	ledger := manager.ServiceConfig{
		Name:     "ledger",
		Address:  a.cfg.LedgerGRPCAddr,
		Insecure: true,
		Timeout:  time.Duration(5 * time.Second),
		Retries:  3,
	}

	audit := manager.ServiceConfig{
		Name:     "audit",
		Address:  a.cfg.AuditGRPCAddr,
		Insecure: true,
		Timeout:  time.Duration(5 * time.Second),
		Retries:  3,
	}

	notify := manager.ServiceConfig{
		Name:     "notification",
		Address:  a.cfg.NotifyGRPCAddr,
		Insecure: true,
		Timeout:  time.Duration(5 * time.Second),
		Retries:  3,
	}

	clients := []manager.ServiceConfig{user, admin, account, transfer, ledger, audit, notify}
	for _, client := range clients {
		c, err := manager.NewConnection(client, a.logger)
		if err != nil {
			return err
		}
		a.manager.Add(client.Name, c)
	}

	return nil
}
