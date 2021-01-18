package config

import (
	"log"

	"github.com/spf13/viper"

	"path"
	"runtime"
)

// Get config path
func getConfigPath() string {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	path := path.Dir(filename)

	return path
}

// Get config
func Config() *viper.Viper {
	config := viper.New()
	config.SetConfigFile(".env")

	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Config file error: %v", err)
	}
	return config
}