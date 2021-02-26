package register

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/request"
)

// RegisteredNode stored network information of a registered node
type RegisteredNode interface {
	GetNodeID() string
	GetAddress() string
	GetRegionCode() string
	GetRootSigningKey() (*fcrcrypto.KeyPair, error)
	GetSigningKey() (*fcrcrypto.KeyPair, error)
	GetNetworkInfoGateway() string
	GetNetworkInfoProvider() string
	GetNetworkInfoClient() string
	GetNetworkInfoAdmin() string
}

// GatewayRegister stores information of a registered gateway
type GatewayRegister struct {
	NodeID              string `json:"nodeId"`
	Address             string `json:"address"`
	RootSigningKey      string `json:"rootSigningKey"`
	SigningKey          string `json:"sigingKey"`
	RegionCode          string `json:"regionCode"`
	NetworkInfoGateway  string `json:"networkInfoGateway"`
	NetworkInfoProvider string `json:"networkInfoProvider"`
	NetworkInfoClient   string `json:"networkInfoClient"`
	NetworkInfoAdmin    string `json:"networkInfoAdmin"`
}

// ProviderRegister stores information of a registered provider
type ProviderRegister struct {
	NodeID             string `json:"nodeId"`
	Address            string `json:"address"`
	RootSigningKey     string `json:"rootSigningKey"`
	SigningKey         string `json:"sigingKey"`
	RegionCode         string `json:"regionCode"`
	NetworkInfoGateway string `json:"networkInfoGateway"`
	NetworkInfoClient  string `json:"networkInfoClient"`
	NetworkInfoAdmin   string `json:"networkInfoAdmin"`
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

// GetNodeID gets the node id
func (r *GatewayRegister) GetNodeID() string {
	return r.NodeID
}

// GetNodeID gets the node id
func (r *ProviderRegister) GetNodeID() string {
	return r.NodeID
}

// GetAddress gets the address
func (r *GatewayRegister) GetAddress() string {
	return r.Address
}

// GetAddress gets the node id
func (r *ProviderRegister) GetAddress() string {
	return r.Address
}

// GetRegionCode gets the region code
func (r *GatewayRegister) GetRegionCode() string {
	return r.RegionCode
}

// GetRegionCode gets the region code
func (r *ProviderRegister) GetRegionCode() string {
	return r.RegionCode
}

// GetNetworkInfoGateway gets the network gateway ap
func (r *GatewayRegister) GetNetworkInfoGateway() string {
	return r.NetworkInfoGateway
}

// GetNetworkInfoGateway gets the network gateway ap
func (r *ProviderRegister) GetNetworkInfoGateway() string {
	return r.NetworkInfoGateway
}

// GetNetworkInfoProvider gets the network provider ap
func (r *GatewayRegister) GetNetworkInfoProvider() string {
	return r.NetworkInfoProvider
}

// GetNetworkInfoProvider gets the network provider ap
func (r *ProviderRegister) GetNetworkInfoProvider() string {
	return ""
}

// GetNetworkInfoClient gets the network client ap
func (r *GatewayRegister) GetNetworkInfoClient() string {
	return r.NetworkInfoClient
}

// GetNetworkInfoClient gets the network client ap
func (r *ProviderRegister) GetNetworkInfoClient() string {
	return r.NetworkInfoClient
}

// GetNetworkInfoAdmin gets the network admin ap
func (r *GatewayRegister) GetNetworkInfoAdmin() string {
	return r.NetworkInfoAdmin
}

// GetNetworkInfoAdmin gets the network admin ap
func (r *ProviderRegister) GetNetworkInfoAdmin() string {
	return r.NetworkInfoAdmin
}

// GetRootSigningKey gets the root signing key
func (r *GatewayRegister) GetRootSigningKey() (*fcrcrypto.KeyPair, error) {
	return fcrcrypto.DecodePublicKey(r.RootSigningKey)
}

// GetRootSigningKey gets the root signing key
func (r *ProviderRegister) GetRootSigningKey() (*fcrcrypto.KeyPair, error) {
	return fcrcrypto.DecodePublicKey(r.RootSigningKey)
}

// GetSigningKey gets the signing key
func (r *GatewayRegister) GetSigningKey() (*fcrcrypto.KeyPair, error) {
	return fcrcrypto.DecodePublicKey(r.SigningKey)
}

// GetSigningKey gets the signing key
func (r *ProviderRegister) GetSigningKey() (*fcrcrypto.KeyPair, error) {
	return fcrcrypto.DecodePublicKey(r.SigningKey)
}

// RegisterGateway to register a gateway
func (r *GatewayRegister) RegisterGateway(registerURL string) error {
	url := registerURL + "/registers/gateway"
	err := request.SendJSON(url, r)
	if err != nil {
		return err
	}
	return nil
}

// RegisterProvider to register a provider
func (r *ProviderRegister) RegisterProvider(registerURL string) error {
	url := registerURL + "/registers/provider"
	err := request.SendJSON(url, r)
	if err != nil {
		return err
	}
	return nil
}
