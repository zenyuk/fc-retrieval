package register

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/request"
)

// GatewayRegister data model
type GatewayRegister struct {
	NodeID              string
	Address             string
	NetworkInfoGateway  string
	NetworkInfoProvider string
	NetworkInfoClient   string
	NetworkInfoAdmin    string
	RegionCode          string
	RootSigningKey      string
	SigingKey           string
}

// ProviderRegister data model
type ProviderRegister struct {
	NodeID             string
	Address            string
	NetworkInfoGateway string
	NetworkInfoClient  string
	NetworkInfoAdmin   string
	RegionCode         string
	RootSigningKey     string
	SigingKey          string
}

// GetRegisteredProviders returns registered providers
func GetRegisteredProviders(registerURL string) ([]ProviderRegister, error) {
	url := registerURL + "/registers/provider"
	providers := []ProviderRegister{}
	err := request.GetJSON(url, &providers)
	if err != nil {
		return providers, err
	}
	return providers, nil
}

// GetRegisteredGateways returns registered gateways
func GetRegisteredGateways(registerURL string) ([]GatewayRegister, error) {
	url := registerURL + "/registers/gateway"
	gateways := []GatewayRegister{}
	err := request.GetJSON(url, &gateways)
	if err != nil {
		return gateways, err
	}
	return gateways, nil
}

// RegisterProvider to register a provider
func RegisterProvider(registerURL string, providerRegister ProviderRegister) error {
	url := registerURL + "/registers/provider"
	err := request.SendJSON(url, providerRegister)
	if err != nil {
		return err
	}
	return nil
}

// RegisterGateway to register a gateway
func RegisterGateway(registerURL string, gatewayRegister GatewayRegister) error {
	url := registerURL + "/registers/gateway"
	err := request.SendJSON(url, gatewayRegister)
	if err != nil {
		return err
	}
	return nil
}
