package config

import (
	"strings"

	commonconfig "github.com/Eucastan/eucastanpay/common/pkg/config"
	"github.com/spf13/viper"
)

func Load() (*Config, error) {

	commonconfig.Init()

	cfg := ToConfig()

	brokers := viper.GetString("KAFKA_BROKERS")
	if brokers != "" {
		cfg.SharedCfg.Kafka.Brokers = strings.Split(brokers, ",")
	}

	if err := commonconfig.RequireString("HTTP_PORT", cfg.HTTPPort); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("USER_GRPC_ADDR", cfg.UserGRPCAddr); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("ADMIN_GRPC_ADDR", cfg.AdminGRPCAddr); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("ACCOUNT_GRPC_ADDR", cfg.AccountGRPCAddr); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("TRANSFER_GRPC_ADDR", cfg.TransferGRPCAddr); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("LEDGER_GRPC_ADDR", cfg.LedgerGRPCAddr); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("AUDIT_GRPC_ADDR", cfg.AuditGRPCAddr); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("NOTIFICATION_GRPC_ADDR", cfg.NotifyGRPCAddr); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireMinLength("JWT_SECRET", cfg.SharedCfg.JWTSecret, 32); err != nil {
		return nil, err
	}

	return cfg, nil
}
