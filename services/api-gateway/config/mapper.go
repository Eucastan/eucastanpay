package config

import (
	"time"

	commonconfig "github.com/Eucastan/eucastanpay/common/pkg/config"
	"github.com/spf13/viper"
)

func ToConfig() *Config {

	timeout := 5 * time.Second

	if v := viper.GetDuration("GRPC_TIMEOUT"); v > 0 {
		timeout = v
	}

	return &Config{
		HTTPPort:    viper.GetString("HTTP_PORT"),
		ServiceName: viper.GetString("SERVICE_NAME"),
		Version:     viper.GetString("VERSION"),

		UserGRPCAddr:     viper.GetString("USER_GRPC_ADDR"),
		AdminGRPCAddr:    viper.GetString("ADMIN_GRPC_ADDR"),
		AccountGRPCAddr:  viper.GetString("ACCOUNT_GRPC_ADDR"),
		TransferGRPCAddr: viper.GetString("TRANSFER_GRPC_ADDR"),
		LedgerGRPCAddr:   viper.GetString("LEDGER_GRPC_ADDR"),
		AuditGRPCAddr:    viper.GetString("AUDIT_GRPC_ADDR"),
		NotifyGRPCAddr:   viper.GetString("NOTIFICATION_GRPC_ADDR"),

		GRPCTimeout:     timeout,
		ShutdownTimeout: 15 * time.Second,

		SharedCfg: commonconfig.SharedCfg{
			JWTSecret: viper.GetString("JWT_SECRET"),
			LogLevel:  viper.GetString("LOG_LEVEL"),
		},
	}
}
