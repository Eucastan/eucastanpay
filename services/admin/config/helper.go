package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func ToCfg() *Config {
	fmt.Println("ToCfg DSN:", viper.GetString("DSN"))

	cfg := &Config{
		Dsn:              viper.GetString("DSN"),
		JWTSecret:        viper.GetString("JWT_SECRET"),
		HTTPPort:         viper.GetString("HTTP_PORT"),
		GRPCPort:         viper.GetString("GRPC_PORT"),
		LedgerGRPCPort:   viper.GetString("LEDGER_GRPC_PORT"),
		AccountGRPCPort:  viper.GetString("ACCOUNT_GRPC_PORT"),
		TransferGRPCPort: viper.GetString("TRANSFER_GRPC_PORT"),
		UserGRPCPort:     viper.GetString("USER_GRPC_PORT"),
		ServiceName:      viper.GetString("SERVICE_NAME"),
		Version:          viper.GetString("VERSION"),
		EmailAPIKey:      viper.GetString("EMAIL_API_KEY"),
		AppEmail:         viper.GetString("APP_EMAIL"),
		FromName:         viper.GetString("FROM_NAME"),
		LogLevel:         viper.GetString("LOG_LEVEL"),
		ShutdownTimeout:  15 * time.Second,
		Redis: Redis{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		Kafka: KafkaConfig{
			Brokers: strings.Split(
				viper.GetString("KAFKA_BROKERS"),
				",",
			),
			Username: viper.GetString("KAFKA_USERNAME"),
			Password: viper.GetString("KAFKA_PASSWORD"),
		},
	}

	fmt.Println("Assigned DSN:", cfg.Dsn)

	return cfg
}
