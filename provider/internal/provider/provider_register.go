package provider

import (
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/request"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
	"github.com/spf13/viper"
)

// Register a provider
func RegisterProvider(conf *viper.Viper) error {
	log.Info("Register provider")
	url := conf.GetString("REGISTER_API_URL") + "/registers/provider"
	reg := register.ProviderRegister{
		Address:        conf.GetString("PROVIDER_ADDRESS"),
		NetworkInfo:    conf.GetString("PROVIDER_NETWORK_INFO"),
		RegionCode:     conf.GetString("PROVIDER_REGION_CODE"),
		RootSigningKey: conf.GetString("PROVIDER_ROOT_SIGNING_KEY"),
		SigingKey:      conf.GetString("PROVIDER_SIGNING_KEY"),
	}
	err := request.SendJSON(url, reg)
	if err != nil {
		log.Error("%+v", err)
		return err
	}
	return nil
}

// Get registered gateways
func GetRegisteredGateways(conf *viper.Viper) ([]register.GatewayRegister, error) {
	url := conf.GetString("REGISTER_API_URL") + "/registers/gateway"
	gateways := []register.GatewayRegister{}
	err := request.GetJSON(url, &gateways)
	if err != nil {
		log.Error("%+v", err)
		return gateways, err
	}
	if len(gateways) == 0 {
		log.Warn("No gateways found")
	}
	return gateways, nil
}
