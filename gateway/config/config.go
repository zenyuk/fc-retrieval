package config

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

func NewConfig() *viper.Viper {
	conf := viper.New()
	conf.AutomaticEnv()
	return conf
}

func parseUint8(value string) uint8 {
	result, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		log.Panic("unable to parse uint8")
	}
	return uint8(result)
}

func Map(conf *viper.Viper) settings.AppSettings {
	return settings.AppSettings{
		BindRestAPI: 				conf.GetString("BIND_REST_API"),
		BindProviderAPI: 		conf.GetString("BIND_PROVIDER_API"),
		BindGatewayAPI: 		conf.GetString("BIND_GATEWAY_API"),
		BindAdminAPI: 			conf.GetString("BIND_ADMIN_API"),
		LogLevel: 					conf.GetString("LOG_LEVEL"),
		LogTarget: 					conf.GetString("LOG_TARGET"),
		LogDir:							conf.GetString("LOG_DIR"),
		LogFile:						conf.GetString("LOG_FILE"),
		LogMaxBackups:			conf.GetInt("LOG_MAX_BACKUPS"),
		LogMaxAge:					conf.GetInt("LOG_MAX_AGE"),
		LogMaxSize:					conf.GetInt("LOG_MAX_SIZE"),
		LogCompress:				conf.GetBool("LOG_COMPRESS"),
		GatewayID: 					conf.GetString("GATEWAY_ID"),
		GatewayPrivKey: 		conf.GetString("GATEWAY_PRIVATE_KEY"),
		GatewayKeyVersion:	conf.GetUint32("GATEWAY_KEY_VERSION"),
		GatewaySigAlg: 			parseUint8(conf.GetString("GATEWAY_SIG_ALG")),
	}
}
