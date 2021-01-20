package logging

import (
	"io"
	"path"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init(conf *viper.Viper) {
	setTimeFormat()
	setLogLevel(conf)
	writer := getLogTarget(conf)
	log.Logger = zerolog.New(writer)
}

func setTimeFormat() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func setLogLevel(conf *viper.Viper) {
	logLevel := conf.GetString("LOG_LEVEL")
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Error().Err(err).Msg("can't parse log level")
		level = zerolog.InfoLevel
		log.Info().Msgf("using level '%v' as default", level)
	}
	zerolog.SetGlobalLevel(level)
}

func getLogTarget(conf *viper.Viper) io.Writer {
	logTarget := conf.GetString("LOG_TARGET")
	switch logTarget {
		case "FILE":   return newLogTargetFile(conf)
		default:       return os.Stdout
	}
}

func newLogTargetFile(conf *viper.Viper) io.Writer {
	logDir := conf.GetString("LOG_DIR")
	if err := os.MkdirAll(logDir, 0744); err != nil {
		log.Error().Err(err).Str("path", logDir).Msg("can't create log directory")
		return nil
	}
	return &lumberjack.Logger{
		Filename:   path.Join(logDir, conf.GetString("LOG_FILE")),
		MaxBackups: conf.GetInt("LOG_MAX_BACKUPS"),
		MaxAge:     conf.GetInt("LOG_MAX_AGE"),
		MaxSize:    conf.GetInt("LOG_MAX_SIZE"),
		Compress:   conf.GetBool("LOG_COMPRESS"),
	}
}

func SetLogLevel(logLevel string) {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Error().Err(err).Str("level", logLevel).Msg("can't parse log level")
	} else {
		zerolog.SetGlobalLevel(level)
	}
}

func Trace(msg string, args ...interface{}) {
	log.Trace().Msgf(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	log.Debug().Msgf(msg, args...)
}

func Info(msg string, args ...interface{}) {
	log.Info().Msgf(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	log.Warn().Msgf(msg, args...)
}

func Error(msg string, args ...interface{}) {
	log.Error().Msgf(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	log.Fatal().Msgf(msg, args...)
}

func Panic(msg string, args ...interface{}) {
	log.Panic().Msgf(msg, args...)
}

/* TODO: Keep to avoid issues */

func ErrorAndPanic(msg string, args ...interface{}) {
	log.Error().Msgf(msg, args...)
}

func Error1(err error) {
	log.Error().Err(err).Msg("Error")
}

/* END TODO */