package bootstrap

import (
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/grpc"
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/discovery"
)

func (a *App) initManager() error {
	mcfg := grpc.ServiceConfig{
		Name:     "account",
		Address:  a.cfg.AccountGRPCPort,
		Insecure: true,
		Timeout:  5 * time.Second,
		Retries:  3,
	}

	m := grpc.NewManager(discovery.NewStaticRegistry(
		map[string]string{
			"account": mcfg.Address,
		},
	))

	conn, err := grpc.NewConnection(mcfg, a.logger)
	if err != nil {
		return err
	}

	m.Add("account", conn)

	a.manager = m
	return nil
}
