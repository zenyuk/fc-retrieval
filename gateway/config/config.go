package config

import (
	"log"
	"strconv"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/spf13/viper"
)

// NewConfig creates a new configuration
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

// Map sets the config for the Gateway. NB: Gateways start without a private key. Private keys are provided by a gateway admin client.
func Map(conf *viper.Viper) settings.AppSettings {
	return settings.AppSettings{
		BindRestAPI:       conf.GetString("BIND_REST_API"),
		BindProviderAPI:   conf.GetString("BIND_PROVIDER_API"),
		BindGatewayAPI:    conf.GetString("BIND_GATEWAY_API"),
		BindAdminAPI:      conf.GetString("BIND_ADMIN_API"),
		LogLevel:          conf.GetString("LOG_LEVEL"),
		LogTarget:         conf.GetString("LOG_TARGET"),
		LogDir:            conf.GetString("LOG_DIR"),
		LogFile:           conf.GetString("LOG_FILE"),
		LogMaxBackups:     conf.GetInt("LOG_MAX_BACKUPS"),
		LogMaxAge:         conf.GetInt("LOG_MAX_AGE"),
		LogMaxSize:        conf.GetInt("LOG_MAX_SIZE"),
		LogCompress:       conf.GetBool("LOG_COMPRESS"),
		GatewayID:         conf.GetString("GATEWAY_ID"),
		GatewayPrivKey:    "",
		GatewayKeyVersion: 0,
		GatewaySigAlg:     parseUint8(conf.GetString("GATEWAY_SIG_ALG")),

		RegisterAPIURL:        conf.GetString("REGISTER_API_URL"),
		GatewayAddress:        conf.GetString("GATEWAY_ADDRESS"),
		GatewayNetworkInfo:    conf.GetString("GATEWAY_NETWORK_INFO"),
		GatewayRegionCode:     conf.GetString("GATEWAY_REGION_CODE"),
		GatewayRootSigningKey: conf.GetString("GATEWAY_ROOT_SIGNING_KEY"),
		GatewaySigningKey:     conf.GetString("GATEWAY_SIGNING_KEY"),

		ClientNetworkInfo:   conf.GetString("CLIENT_NETWORK_INFO"),
		ProviderNetworkInfo: conf.GetString("PROVIDER_NETWORK_INFO"),
		AdminNetworkInfo:    conf.GetString("ADMIN_NETWORK_INFO"),
	}
}
