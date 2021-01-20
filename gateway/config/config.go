package config

import (
	"log"

	"github.com/spf13/viper"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

func NewConfig() settings.AppSettings {
	viper := viper.New()
  viper.SetConfigName("settings")
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/gateway/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("Fatal error config file")
	}
	var config settings.AppSettings
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
  return config
}
