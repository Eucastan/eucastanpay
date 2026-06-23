package config

import (
	commonconfig "github.com/Eucastan/eucastanpay/common/pkg/config"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Dsn         string `mapstructure:"DSN"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	HTTPPort    string `mapstructure:"HTTP_PORT"`
	GRPCPort    string `mapstructure:"GRPC_PORT"`
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Version     string `mapstructure:"VERSION"`
	Redis       Redis  `mapstructure:",squash"`
	EmailAPIKey string `mapstructure:"EMAIL_API_KEY"`
	AppEmail    string `mapstructure:"APP_EMAIL"`
	FromName    string `mapstructure:"FROM_NAME"`
	Kafka       KafkaConfig
	LogLevel    string `mapstructure:"LOG_LEVEL"`
}

type Redis struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"KAFKA_BROKERS"`
}

func Load() (*Config, error) {
	commonconfig.Init()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

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

	if err := commonconfig.RequireMinLength("JWT_SECRET", cfg.GRPCPort, 32); err != nil {
		return nil, err
	}

	return &cfg, nil
}
