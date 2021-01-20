package logger

import (
  "github.com/rs/zerolog"
  "github.com/spf13/viper"
)

func InitLogger(conf *viper.Viper) {
  setTimeFormat()
  setLogLevel(conf)
}

func setTimeFormat() {
  zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func setLogLevel(conf *viper.Viper) {
  logLevel := conf.GetString("LOGGING_LEVEL")
  switch logLevel {
    case "debug": zerolog.SetGlobalLevel(zerolog.DebugLevel)
    case "trace": zerolog.SetGlobalLevel(zerolog.TraceLevel)
    case "warn":  zerolog.SetGlobalLevel(zerolog.WarnLevel)
    case "error": zerolog.SetGlobalLevel(zerolog.ErrorLevel)
    case "fatal": zerolog.SetGlobalLevel(zerolog.FatalLevel)
    case "panic": zerolog.SetGlobalLevel(zerolog.PanicLevel)
    default:      zerolog.SetGlobalLevel(zerolog.InfoLevel)
  }
}