package config

import (
	"errors"
	"github.com/spf13/viper"
)

type Config struct {
	Dsn             string `mapstructure:"DSN"`
	JWTSecret       string `mapstructure:"JWT_SECRET"`
	HTTPPort        string `mapstructure:"HTTP_PORT"`
	GRPCPort        string `mapstructure:"GRPC_PORT"`
	LedgerGRPCPort  string `mapstructure:"LEDGER_GRPC_PORT"`
	AccountGRPCPort string `mapstructure:"ACCOUNT_GRPC_PORT"`
	UserGRPCPort    string `mapstructure:"USER_GRPC_PORT"`
	ServiceName     string `mapstructure:"SERVICE_NAME"`
	Version         string `mapstructure:"VERSION"`
	Redis           Redis
	EmailAPIKey     string `mapstructure:"EMAIL_API_KEY"`
	AppEmail        string `mapstructure:"APP_EMAIL"`
	FromName        string `mapstructure:"FROM_NAME"`
	Kafka           KafkaConfig
	LogLevel        string `mapstructure:"LOG_LEVEL"`
}

type Redis struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type KafkaConfig struct{ Brokers []string }

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if len(cfg.JWTSecret) < 32 {
		return nil, errors.New("invalid JWT Secret")
	}

	return &cfg, nil
}
