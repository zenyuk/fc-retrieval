package config

import (
  "github.com/spf13/viper"
)

func Config() *viper.Viper {
  config := viper.New()
  config.AutomaticEnv()
  return config
}