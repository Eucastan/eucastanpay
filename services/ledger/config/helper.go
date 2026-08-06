package config

import (
	"fmt"
	"strings"
	"time"

	commonconfig "github.com/Eucastan/eucastanpay/common/pkg/config"
	"github.com/spf13/viper"
)

func ToCfg() *Config {
	fmt.Println("ToCfg DSN:", viper.GetString("DSN"))

	cfg := &Config{
		HTTPPort:        viper.GetString("HTTP_PORT"),
		GRPCPort:        viper.GetString("GRPC_PORT"),
		AccountGRPCPort: viper.GetString("ACCOUNT_GRPC_PORT"),
		ServiceName:     viper.GetString("SERVICE_NAME"),
		Version:         viper.GetString("VERSION"),
		EmailAPIKey:     viper.GetString("EMAIL_API_KEY"),
		AppEmail:        viper.GetString("APP_EMAIL"),
		FromName:        viper.GetString("FROM_NAME"),
		ShutdownTimeout: 15 * time.Second,
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
				GroupID:  viper.GetString("KAFKA_GROUP_ID"),
				Username: viper.GetString("KAFKA_USERNAME"),
				Password: viper.GetString("KAFKA_PASSWORD"),
				CaCert:   viper.GetString("KAFKA_CA_CERT"),
				TLS:      true,
				SASL:     true,
			},
			LogLevel: viper.GetString("LOG_LEVEL"),
		},
	}

	fmt.Println("Assigned DSN:", cfg.SharedCfg.Dsn)
	fmt.Println("Assigned Account Service Address:", cfg.AccountGRPCPort)

	return cfg
}
