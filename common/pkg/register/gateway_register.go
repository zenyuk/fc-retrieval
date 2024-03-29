/*
Package register - provides methods for FileCoin Secondary Retrieval from Retrieval Register perspective.

Retrieval Register is a central node, holding information about Retrieval Gateways and Retrieval Providers.
*/
package register

import (
  "github.com/ConsenSys/fc-retrieval/common/pkg/fcrcrypto"
)

// GatewayRegister stores information of a registered gateway
type GatewayRegister struct {
  NodeID              string `json:"nodeId"`
  Address             string `json:"address"`
  RootSigningKey      string `json:"rootSigningKey"`
  SigningKey          string `json:"signingKey"`
  RegionCode          string `json:"regionCode"`
  NetworkInfoGateway  string `json:"networkInfoGateway"`
  NetworkInfoProvider string `json:"networkInfoProvider"`
  NetworkInfoClient   string `json:"networkInfoClient"`
  NetworkInfoAdmin    string `json:"networkInfoAdmin"`
}

type GatewayRegistrar interface {
  GetNodeID() string
  GetAddress() string
  GetRegionCode() string
  GetNetworkInfoGateway() string
  GetNetworkInfoProvider() string
  GetNetworkInfoClient() string
  GetNetworkInfoAdmin() string
  GetRootSigningKey() (*fcrcrypto.KeyPair, error)
  GetSigningKey() (*fcrcrypto.KeyPair, error)
  Serialize() GatewayRegister
}

func NewGatewayRegister(
  nodeID              string,
  address             string,
  rootSigningKey      string,
  signingKey          string,
  regionCode          string,
  networkInfoGateway  string,
  networkInfoProvider string,
  networkInfoClient   string,
  networkInfoAdmin    string,
  ) GatewayRegistrar {
  return &GatewayRegister {
    NodeID: nodeID,
    Address: address,
    RootSigningKey: rootSigningKey,
    SigningKey: signingKey,
    RegionCode: regionCode,
    NetworkInfoGateway: networkInfoGateway,
    NetworkInfoProvider: networkInfoProvider,
    NetworkInfoClient: networkInfoClient,
    NetworkInfoAdmin: networkInfoAdmin,
  }
}

// GetNodeID gets the node id
func (r *GatewayRegister) GetNodeID() string {
  return r.NodeID
}

// GetAddress gets the address
func (r *GatewayRegister) GetAddress() string {
  return r.Address
}

// GetRegionCode gets the region code
func (r *GatewayRegister) GetRegionCode() string {
  return r.RegionCode
}

// GetNetworkInfoGateway gets the network gateway ap
func (r *GatewayRegister) GetNetworkInfoGateway() string {
  return r.NetworkInfoGateway
}

// GetNetworkInfoProvider gets the network provider ap
func (r *GatewayRegister) GetNetworkInfoProvider() string {
  return r.NetworkInfoProvider
}

// GetNetworkInfoClient gets the network client ap
func (r *GatewayRegister) GetNetworkInfoClient() string {
  return r.NetworkInfoClient
}

// GetNetworkInfoAdmin gets the network admin ap
func (r *GatewayRegister) GetNetworkInfoAdmin() string {
  return r.NetworkInfoAdmin
}

// GetRootSigningKey gets the root signing key
func (r *GatewayRegister) GetRootSigningKey() (*fcrcrypto.KeyPair, error) {
  return fcrcrypto.DecodePublicKey(r.RootSigningKey)
}

// GetSigningKey gets the signing key
func (r *GatewayRegister) GetSigningKey() (*fcrcrypto.KeyPair, error) {
  return fcrcrypto.DecodePublicKey(r.SigningKey)
}

func (r *GatewayRegister) Serialize() GatewayRegister {
  return *r
}
