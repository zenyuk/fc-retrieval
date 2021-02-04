package register

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/request"
)

// GatewayRegister data model
type GatewayRegister struct {
	NodeID              string
	Address             string
	NetworkGatewayInfo  string
	NetworkProviderInfo string
	NetworkClientInfo   string
	NetworkAdminInfo    string
	RegionCode          string
	RootSigningKey      string
	SigingKey           string
}

// ProviderRegister data model
type ProviderRegister struct {
	NodeID         string
	Address        string
	NetworkInfo    string
	RegionCode     string
	RootSigningKey string
	SigingKey      string
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
