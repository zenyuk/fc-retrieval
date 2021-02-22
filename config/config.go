package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewConfig(file string) *viper.Viper {
	config := viper.New()
	config.AddConfigPath("../../")

	config.SetConfigName(file)
	config.SetConfigType("env")
	err := config.ReadInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return config
}
