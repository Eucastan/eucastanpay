package config

import (
	"strings"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigFile(".env")

	// Optional
	_ = viper.ReadInConfig()

	viper.SetEnvKeyReplacer(
		strings.NewReplacer(".", "_"),
	)

	viper.AutomaticEnv()
}
