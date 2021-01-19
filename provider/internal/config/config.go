package config

import (
  "github.com/spf13/viper"
  "go.uber.org/fx"
)

func NewConfig() *viper.Viper {
  config := viper.New()
  config.AutomaticEnv()
  return config
}

var Module = fx.Options(
  fx.Provide(
    NewConfig,
  ),
)
