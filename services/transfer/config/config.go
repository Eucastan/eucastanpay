package config

import (
	"strings"
	"time"

	commonconfig "github.com/Eucastan/eucastanpay/common/pkg/config"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPPort        string `mapstructure:"HTTP_PORT"`
	GRPCADDR        string `mapstructure:"GRPC_ADDR"`
	LedgerGRPCADDR  string `mapstructure:"LEDGER_GRPC_ADDR"`
	AccountGRPCADDR string `mapstructure:"ACCOUNT_GRPC_ADDR"`
	UserGRPCADDR    string `mapstructure:"USER_GRPC_ADDR"`
	ServiceName     string `mapstructure:"SERVICE_NAME"`
	Version         string `mapstructure:"VERSION"`
	EmailAPIKey     string `mapstructure:"EMAIL_API_KEY"`
	AppEmail        string `mapstructure:"APP_EMAIL"`
	FromName        string `mapstructure:"FROM_NAME"`
	LogLevel        string `mapstructure:"LOG_LEVEL"`
	ShutdownTimeout time.Duration
	SharedCfg       commonconfig.SharedCfg
}

func Load() (*Config, error) {
	commonconfig.Init()

	cfg := ToCfg()

	brokers := viper.GetString("KAFKA_BROKERS")
	if brokers != "" {
		cfg.SharedCfg.Kafka.Brokers = strings.Split(brokers, ",")
	}

	if err := commonconfig.RequireString("DSN", cfg.SharedCfg.Dsn); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("HTTP_PORT", cfg.HTTPPort); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("GRPC_PORT", cfg.GRPCADDR); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("ACCOUNT_GRPC_PORT", cfg.AccountGRPCADDR); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("LEDGER_GRPC_PORT", cfg.LedgerGRPCADDR); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("USER_GRPC_PORT", cfg.UserGRPCADDR); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireMinLength("JWT_SECRET", cfg.SharedCfg.JWTSecret, 32); err != nil {
		return nil, err
	}

	return cfg, nil
}
