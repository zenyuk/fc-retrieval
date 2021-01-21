package config

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

func NewConfig() settings.AppSettings {
	config := viper.New()
	config.AutomaticEnv()
	return settings.AppSettings{
		BindRestAPI: 				config.GetString("BIND_REST_API"),
		BindProviderAPI: 		config.GetString("BIND_PROVIDER_API"),
		BindGatewayAPI: 		config.GetString("BIND_GATEWAY_API"),
		BindAdminAPI: 			config.GetString("BIND_ADMIN_API"),
		LogLevel: 					config.GetString("LOG_LEVEL"),
		LogTarget: 					config.GetString("LOG_TARGET"),
		GatewayID: 					config.GetString("GATEWAY_ID"),
		GatewayPrivKey: 		config.GetString("GATEWAY_PRIVATE_KEY"),
		GatewayKeyVersion:	config.GetUint32("GATEWAY_KEY_VERSION"),
		GatewaySigAlg: 			parseUint8(config.GetString("GATEWAY_SIG_ALG")),
	}
}

func parseUint8(value string) uint8 {
	result, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		log.Panic("unable to parse uint8")
	}
	return uint8(result)
}
