package clients

import (
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/grpc"
	"github.com/Eucastan/eucastanpay/services/transfer/config"
)

func Init(cfg *config.Config) grpc.ServiceConfig {
	return grpc.ServiceConfig{
		Name:     "account",
		Address:  cfg.AccountGRPCADDR,
		Insecure: true,
		Timeout:  5 * time.Second,
		Retries:  3,
	}
}
