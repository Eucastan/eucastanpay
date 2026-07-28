package bootstrap

import (
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/grpc"
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/discovery"
)

func (a *App) initManager() error {
	account := grpc.ServiceConfig{
		Name:     "account",
		Address:  a.cfg.AccountGRPCPort,
		Insecure: true,
		Timeout:  5 * time.Second,
		Retries:  3,
	}

	m := grpc.NewManager(discovery.NewStaticRegistry(
		map[string]string{
			"account": account.Address,
		},
	))

	accountConn, err := grpc.NewConnection(account, a.logger)
	if err != nil {
		return err
	}

	m.Add(account.Name, accountConn)

	a.manager = m
	return nil
}
