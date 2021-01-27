package config

import (
	"github.com/spf13/viper"
)

// NewConfig provides a configuration environment
func NewConfig() *viper.Viper {
	config := viper.New()
	config.AutomaticEnv()
	return config
}
