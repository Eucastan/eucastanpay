package clients

import (
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/grpc/clients"
	"github.com/Eucastan/eucastanpay/services/transfer/config"
	"github.com/sirupsen/logrus"
)

func Init(cfg *config.Config, log *logrus.Logger) (*clients.Clients, error) {
	clientCfg := clients.Config{
		UserServiceAddr:    cfg.UserGRPCPort,
		AccountServiceAddr: cfg.AccountGRPCPort,
		LedgerServiceAddr:  cfg.LedgerGRPCPort,
		Timeout:            5 * time.Second,
		MaxRetries:         3,
		Insecure:           true,
	}

	return clients.NewClients(clientCfg, log)
}
