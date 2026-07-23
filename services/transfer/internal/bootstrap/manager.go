package bootstrap

import (
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/grpc"
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/discovery"
)

func (a *App) initManager() error {
	account := grpc.ServiceConfig{
		Name:     "account",
		Address:  a.cfg.AccountGRPCADDR,
		Insecure: true,
		Timeout:  5 * time.Second,
		Retries:  3,
	}

	ledger := grpc.ServiceConfig{
		Name:     "ledger",
		Address:  a.cfg.LedgerGRPCADDR,
		Insecure: true,
		Timeout:  5 * time.Second,
		Retries:  3,
	}

	m := grpc.NewManager(discovery.NewStaticRegistry(
		map[string]string{
			"account": account.Address,
			"ledger":  ledger.Address,
		},
	))

	accountConn, err := grpc.NewConnection(account, a.logger)
	if err != nil {
		return err
	}

	ledgerConn, err := grpc.NewConnection(ledger, a.logger)
	if err != nil {
		accountConn.Close()
		return err
	}

	m.Add(account.Name, accountConn)
	m.Add(ledger.Name, ledgerConn)
	a.manager = m

	return nil
}
