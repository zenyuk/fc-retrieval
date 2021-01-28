package register

import (
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/request"
	"github.com/spf13/viper"
)

// Register data model
type Register struct {
	NodeID				 string
	Address        string
	NetworkInfo    string
	RegionCode     string
	RootSigningKey string
	SigingKey      string
}

// Register a provider
func RegisterProvider(conf *viper.Viper) error {
	log.Info("Register provider")
	url := conf.GetString("REGISTER_API_URL") + "/registers/provider"
	providerReg := Register{
		Address:        conf.GetString("PROVIDER_ADDRESS"),
		NetworkInfo:    conf.GetString("PROVIDER_NETWORK_INFO"),
		RegionCode:     conf.GetString("PROVIDER_REGION_CODE"),
		RootSigningKey: conf.GetString("PROVIDER_ROOT_SIGNING_KEY"),
		SigingKey:      conf.GetString("PROVIDER_SIGNING_KEY"),
	}
	err := request.SendJSON(url, providerReg)
	if err != nil {
		log.Error("%+v", err)
		return err
	}
	return nil
}

// Get registered gateways
func GetRegisteredGateways(conf *viper.Viper) ([]Register, error) {
	url := conf.GetString("REGISTER_API_URL") + "/registers/gateway"
	gateways := []Register{}
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
