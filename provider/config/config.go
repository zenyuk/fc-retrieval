package config

import (
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/ConsenSys/fc-retrieval-provider/internal/util/settings"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// NewConfig creates a new configuration
func NewConfig() *viper.Viper {
	conf := viper.New()
	conf.AutomaticEnv()
	defineFlags(conf)
	bindFlags(conf)
	setValues(conf)
	return conf
}

func parseUint8(value string) uint8 {
	result, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		log.Panic("unable to parse uint8")
	}
	return uint8(result)
}

// Map sets the config for the Provider. NB: Providers start without a private key. Private keys are provided by a provider admin client.
func Map(conf *viper.Viper) settings.AppSettings {
	registerRefreshDuration, err := time.ParseDuration(conf.GetString("REGISTER_REFRESH_DURATION"))
	if err != nil {
		registerRefreshDuration = settings.DefaultRegisterRefreshDuration
	}
	tcpInactivityTimeout, err := time.ParseDuration(conf.GetString("TCP_INACTIVITY_TIMEOUT"))
	if err != nil {
		tcpInactivityTimeout = settings.DefaultTCPInactivityTimeout
	}
	tcpLongInactivityTimeout, err := time.ParseDuration(conf.GetString("TCP_LONG_INACTIVITY_TIMEOUT"))
	if err != nil {
		tcpLongInactivityTimeout = settings.DefaultLongTCPInactivityTimeout
	}
	return settings.AppSettings{
		BindRestAPI:    conf.GetString("BIND_REST_API"),
		BindGatewayAPI: conf.GetString("BIND_GATEWAY_API"),
		BindAdminAPI:   conf.GetString("BIND_ADMIN_API"),

		LogLevel:      conf.GetString("LOG_LEVEL"),
		LogTarget:     conf.GetString("LOG_TARGET"),
		LogDir:        conf.GetString("LOG_DIR"),
		LogFile:       conf.GetString("LOG_FILE"),
		LogMaxBackups: conf.GetInt("LOG_MAX_BACKUPS"),
		LogMaxAge:     conf.GetInt("LOG_MAX_AGE"),
		LogMaxSize:    conf.GetInt("LOG_MAX_SIZE"),
		LogCompress:   conf.GetBool("LOG_COMPRESS"),

		ProviderID:     conf.GetString("PROVIDER_ID"),
		ProviderSigAlg: parseUint8(conf.GetString("PROVIDER_SIG_ALG")),

		RegisterAPIURL:         	conf.GetString("REGISTER_API_URL"),
		RegisterRefreshDuration:	registerRefreshDuration,
		
		ProviderAddress:        conf.GetString("PROVIDER_ADDRESS"),
		ProviderRootSigningKey: conf.GetString("PROVIDER_ROOT_SIGNING_KEY"),
		ProviderSigningKey:     conf.GetString("PROVIDER_SIGNING_KEY"),
		ProviderRegionCode:     conf.GetString("PROVIDER_REGION_CODE"),
		NetworkInfoClient:      conf.GetString("IP") + ":" + conf.GetString("NETWORK_CLIENT_INFO"),
		NetworkInfoGateway:     conf.GetString("IP") + ":" + conf.GetString("NETWORK_GATEWAY_INFO"),
		NetworkInfoAdmin:       conf.GetString("IP") + ":" + conf.GetString("NETWORK_ADMIN_INFO"),

		TCPInactivityTimeout: 		tcpInactivityTimeout,
		TCPLongInactivityTimeout: tcpLongInactivityTimeout,
	}
}

func defineFlags(conf *viper.Viper) {
	flag.String("host", "0.0.0.0", "help message for host")
	flag.String("ip", "127.0.0.1", "help message for ip")
}

func bindFlags(conf *viper.Viper) {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	conf.BindPFlags(pflag.CommandLine)
}

func setValues(conf *viper.Viper) {
	conf.Set("IP", conf.GetString("ip"))
}
