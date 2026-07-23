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
		HTTPPort:    viper.GetString("HTTP_PORT"),
		GRPCPort:    viper.GetString("GRPC_PORT"),
		ServiceName: viper.GetString("SERVICE_NAME"),
		Version:     viper.GetString("VERSION"),
		EmailAPIKey: viper.GetString("EMAIL_API_KEY"),
		AppEmail:    viper.GetString("APP_EMAIL"),
		FromName:    viper.GetString("FROM_NAME"),
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
		ShutdownTimeout: 15 * time.Second,
	}

	fmt.Println("Assigned DSN:", cfg.SharedCfg.Dsn)

	return cfg
}
