package logging

import (
  log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
  "github.com/spf13/viper"
)

func InitLogger(conf *viper.Viper) {
	log.Init()
  setLogLevel(conf)
}

func setLogLevel(conf *viper.Viper) {
	logLevel := conf.GetString("LOG_LEVEL")
	log.SetLogLevel(logLevel)
}