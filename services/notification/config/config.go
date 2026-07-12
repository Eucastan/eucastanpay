package config

import (
	commonconfig "github.com/Eucastan/eucastanpay/common/pkg/config"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Dsn         string
	JWTSecret   string
	HTTPPort    string
	GRPCPort    string
	ServiceName string
	Version     string
	Redis       Redis `mapstructure:",squash"`
	EmailAPIKey string
	AppEmail    string
	FromName    string
	Kafka       KafkaConfig
	LogLevel    string
	GinMode     string
}

type Redis struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type KafkaConfig struct {
	Brokers  []string `mapstructure:"KAFKA_BROKERS"`
	Username string   `mapstructure:"KAFKA_USERNAME"`
	Password string   `mapstructure:"KAFKA_PASSWORD"`
}

func Load() (*Config, error) {
	commonconfig.Init()

	cfg := ToCfg()

	brokers := viper.GetString("KAFKA_BROKERS")
	if brokers != "" {
		cfg.Kafka.Brokers = strings.Split(brokers, ",")
	}

	if err := commonconfig.RequireString("DSN", cfg.Dsn); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("HTTP_PORT", cfg.HTTPPort); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireString("GRPC_PORT", cfg.GRPCPort); err != nil {
		return nil, err
	}

	if err := commonconfig.RequireMinLength("JWT_SECRET", cfg.JWTSecret, 32); err != nil {
		return nil, err
	}

	return cfg, nil
}
