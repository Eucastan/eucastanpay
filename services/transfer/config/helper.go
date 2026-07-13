package config

import (
	"fmt"
	commonconfig "github.com/Eucastan/eucastanpay/common/pkg/config"
	"github.com/spf13/viper"
	"strings"
)

func ToCfg() *Config {
	fmt.Println("ToCfg DSN:", viper.GetString("DSN"))

	cfg := &Config{
		HTTPPort:        viper.GetString("HTTP_PORT"),
		GRPCADDR:        viper.GetString("GRPC_ADDR"),
		LedgerGRPCADDR:  viper.GetString("LEDGER_GRPC_ADDR"),
		AccountGRPCADDR: viper.GetString("ACCOUNT_GRPC_ADDR"),
		UserGRPCADDR:    viper.GetString("USER_GRPC_ADDR"),
		ServiceName:     viper.GetString("SERVICE_NAME"),
		Version:         viper.GetString("VERSION"),
		EmailAPIKey:     viper.GetString("EMAIL_API_KEY"),
		AppEmail:        viper.GetString("APP_EMAIL"),
		FromName:        viper.GetString("FROM_NAME"),
		SharedCfg: commonconfig.SharedCfg{
			Dsn:       viper.GetString("DSN"),
			Schema:    viper.GetString("SCHEMA"),
			JWTSecret: viper.GetString("JWT_SECRET"),
			Redis: commonconfig.Redis{
				Addr:     viper.GetString("REDIS_ADDR"),
				Password: viper.GetString("REDIS_PASSWORD"),
				DB:       viper.GetInt("REDIS_DB"),
			},
			Kafka: commonconfig.KafkaConfig{
				Brokers: strings.Split(
					viper.GetString("KAFKA_BROKERS"),
					",",
				),
				Username: viper.GetString("KAFKA_USERNAME"),
				Password: viper.GetString("KAFKA_PASSWORD"),
			},
			LogLevel: viper.GetString("LOG_LEVEL"),
		},
	}

	fmt.Println("Assigned DSN:", cfg.SharedCfg.Dsn)

	return cfg
}
