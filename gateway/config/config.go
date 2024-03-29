/*
Package config - combines operations used to setup parameters for Gateway node in FileCoin network
*/
package config

import (
  "flag"
  "fmt"
  "math/big"
  "time"

  "github.com/spf13/pflag"
  "github.com/spf13/viper"

  "github.com/ConsenSys/fc-retrieval/common/pkg/logging"
  "github.com/ConsenSys/fc-retrieval/gateway/internal/util/settings"
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

// Map sets the config for the Gateway. NB: Gateways start without a private key. Private keys are provided by a gateway admin client.
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

	defaultSearchPrice := new(big.Int)
	_, err = fmt.Sscan(conf.GetString("SEARCH_PRICE"), defaultSearchPrice)
	if err != nil {
		// defaultSearchPrice is the default search price "0.001".
		defaultSearchPrice = big.NewInt(1_000_000_000_000_000)
	}

	defaultOfferPrice := new(big.Int)
	_, err = fmt.Sscan(conf.GetString("OFFER_PRICE"), defaultOfferPrice)
	if err != nil {
		// defaultOfferPrice is the default offer price "0.001".
		defaultOfferPrice = big.NewInt(1_000_000_000_000_000)
	}

	defaultTopUpAmount := new(big.Int)
	_, err = fmt.Sscan(conf.GetString("TOPUP_AMOUNT"), defaultTopUpAmount)
	if err != nil {
		// defaultTopUpAmount is the default top up amount "0.1".
		defaultTopUpAmount = big.NewInt(100_000_000_000_000_000)
	}

	return settings.AppSettings{
		BindRestAPI:     conf.GetString("BIND_REST_API"),
		BindProviderAPI: conf.GetString("BIND_PROVIDER_API"),
		BindGatewayAPI:  conf.GetString("BIND_GATEWAY_API"),
		BindAdminAPI:    conf.GetString("BIND_ADMIN_API"),
		LogLevel:        conf.GetString("LOG_LEVEL"),
		LogTarget:       conf.GetString("LOG_TARGET"),
		LogDir:          conf.GetString("LOG_DIR"),
		LogFile:         conf.GetString("LOG_FILE"),
		LogMaxBackups:   conf.GetInt("LOG_MAX_BACKUPS"),
		LogMaxAge:       conf.GetInt("LOG_MAX_AGE"),
		LogMaxSize:      conf.GetInt("LOG_MAX_SIZE"),
		LogCompress:     conf.GetBool("LOG_COMPRESS"),
		GatewayID:       conf.GetString("GATEWAY_ID"),

		RegisterAPIURL:          conf.GetString("REGISTER_API_URL"),
		RegisterRefreshDuration: registerRefreshDuration,

		GatewayAddress:        conf.GetString("GATEWAY_ADDRESS"),
		NetworkInfoGateway:    conf.GetString("IP") + ":" + conf.GetString("BIND_GATEWAY_API"),
		GatewayRegionCode:     conf.GetString("GATEWAY_REGION_CODE"),
		GatewayRootSigningKey: conf.GetString("GATEWAY_ROOT_SIGNING_KEY"),
		GatewaySigningKey:     conf.GetString("GATEWAY_SIGNING_KEY"),

		NetworkInfoClient:   conf.GetString("IP") + ":" + conf.GetString("BIND_REST_API"),
		NetworkInfoProvider: conf.GetString("IP") + ":" + conf.GetString("BIND_PROVIDER_API"),
		NetworkInfoAdmin:    conf.GetString("IP") + ":" + conf.GetString("BIND_ADMIN_API"),

		TCPInactivityTimeout:     tcpInactivityTimeout,
		TCPLongInactivityTimeout: tcpLongInactivityTimeout,

		SearchPrice: defaultSearchPrice,
		OfferPrice:  defaultOfferPrice,
		TopupAmount: defaultTopUpAmount,
	}
}

func defineFlags(conf *viper.Viper) {
	flag.String("host", "0.0.0.0", "help message for host")
	flag.String("ip", "127.0.0.1", "help message for ip")
}

func bindFlags(conf *viper.Viper) {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := conf.BindPFlags(pflag.CommandLine); err != nil {
	  logging.Error("can't bind a command line flag")
  }
}

func setValues(conf *viper.Viper) {
	conf.Set("IP", conf.GetString("ip"))
}
