package config

import (
	"errors"
	"fmt"
)

type SharedCfg struct {
	Dsn       string `mapstructure:"DSN"`
	Schema    string `mapstructure:"SCHEMA"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
	Redis     Redis  `mapstructure:",squash"`
	Kafka     KafkaConfig
	LogLevel  string `mapstructure:"LOG_LEVEL"`
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

func ValidateSchema(schema string) error {
	switch schema {
	case "identity",
		"account",
		"ledger",
		"transfer",
		"audit",
		"payment",
		"admin":
		return nil
	default:
		return fmt.Errorf("invalid schema: %s", schema)
	}
}

func RequireString(name, value string) error {
	if value == "" {
		return errors.New(name + " is required")
	}

	return nil
}

func RequireMinLength(name, value string, min int) error {
	if len(value) < min {
		return errors.New(name + " is invalid")
	}

	return nil
}
