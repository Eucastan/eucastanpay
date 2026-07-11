package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func ToCfg() *Config {
	fmt.Println("ToCfg DSN:", viper.GetString("DSN"))

	cfg := &Config{
		Dsn:         viper.GetString("DSN"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		HTTPPort:    viper.GetString("HTTP_PORT"),
		GRPCPort:    viper.GetString("GRPC_PORT"),
		ServiceName: viper.GetString("SERVICE_NAME"),
		Version:     viper.GetString("VERSION"),
		EmailAPIKey: viper.GetString("EMAIL_API_KEY"),
		AppEmail:    viper.GetString("APP_EMAIL"),
		FromName:    viper.GetString("FROM_NAME"),
		LogLevel:    viper.GetString("LOG_LEVEL"),
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
