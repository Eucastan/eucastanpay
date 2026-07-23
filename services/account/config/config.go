package config

import (
	"strings"
	"time"

	commonconfig "github.com/Eucastan/eucastanpay/common/pkg/config"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPPort        string `mapstructure:"HTTP_PORT"`
	GRPCPort        string `mapstructure:"GRPC_PORT"`
	ServiceName     string `mapstructure:"SERVICE_NAME"`
	Version         string `mapstructure:"VERSION"`
	EmailAPIKey     string `mapstructure:"EMAIL_API_KEY"`
	AppEmail        string `mapstructure:"APP_EMAIL"`
	FromName        string `mapstructure:"FROM_NAME"`
	ShutdownTimeout time.Duration
	SharedCfg       commonconfig.SharedCfg
}

func Load() (*Config, error) {
	commonconfig.Init()
	println("DSN =", viper.GetString("DSN"))
	println("HTTP_PORT =", viper.GetString("HTTP_PORT"))

	cfg := ToCfg()

	println("cfg.Dsn =", cfg.SharedCfg.Dsn)

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

	if err := commonconfig.RequireString("GRPC_PORT", cfg.GRPCPort); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireMinLength("JWT_SECRET", cfg.SharedCfg.JWTSecret, 32); err != nil {
		return nil, err
	}

	return cfg, nil
}
