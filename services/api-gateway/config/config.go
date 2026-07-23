package config

import (
	"time"

	commonconfig "github.com/Eucastan/eucastanpay/common/pkg/config"
)

type Config struct {
	HTTPPort    string `mapstructure:"HTTP_PORT"`
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Version     string `mapstructure:"VERSION"`

	UserGRPCAddr     string `mapstructure:"USER_GRPC_ADDR"`
	AdminGRPCAddr    string `mapstructure:"ADMIN_GRPC_ADDR"`
	AccountGRPCAddr  string `mapstructure:"ACCOUNT_GRPC_ADDR"`
	TransferGRPCAddr string `mapstructure:"TRANSFER_GRPC_ADDR"`
	LedgerGRPCAddr   string `mapstructure:"LEDGER_GRPC_ADDR"`
	AuditGRPCAddr    string `mapstructure:"AUDIT_GRPC_ADDR"`
	NotifyGRPCAddr   string `mapstructure:"NOTIFICATION_GRPC_ADDR"`

	GRPCTimeout     time.Duration
	ShutdownTimeout time.Duration

	SharedCfg commonconfig.SharedCfg
}
