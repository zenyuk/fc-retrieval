package adminapi

import (
  "github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
  "github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
  "github.com/ConsenSys/fc-retrieval-common/pkg/register"
  "github.com/ConsenSys/fc-retrieval-common/pkg/request"
)

type Admin struct {
  httpCommunicator request.HttpCommunications
}

type AdminApi interface {
  RequestForceRefresh(
    gatewayRegistrar register.GatewayRegistrar,
    signingPrivkey *fcrcrypto.KeyPair,
    signingPrivKeyVer *fcrcrypto.KeyVersion,
  ) error

  RequestGetReputation(
    gatewayRegistrar register.GatewayRegistrar,
    clientID *nodeid.NodeID,
    signingPrivkey *fcrcrypto.KeyPair,
    signingPrivKeyVer *fcrcrypto.KeyVersion,
  ) (int64, error)

  RequestInitialiseKey(
    gatewayRegistrar register.GatewayRegistrar,
    gatewayPrivKey *fcrcrypto.KeyPair,
    gatewayPrivKeyVer *fcrcrypto.KeyVersion,
    signingPrivkey *fcrcrypto.KeyPair,
    signingPrivKeyVer *fcrcrypto.KeyVersion,
  ) error

  RequestInitialiseKeyV2(
    gatewayRegistrar register.GatewayRegistrar,
    gatewayPrivKey *fcrcrypto.KeyPair,
    gatewayPrivKeyVer *fcrcrypto.KeyVersion,
    signingPrivkey *fcrcrypto.KeyPair,
    signingPrivKeyVer *fcrcrypto.KeyVersion,
    lotusWalletPrivateKey string,
    lotusAP string,
    lotusAuthToken string,
  ) error

  RequestListDHTOffer(
    gatewayRegistrar register.GatewayRegistrar,
    signingPrivkey *fcrcrypto.KeyPair,
    signingPrivKeyVer *fcrcrypto.KeyVersion,
  ) error

  SetGroupCIDOfferSupportedForProviders(
    gatewayRegistrar register.GatewayRegistrar,
    providerIDs []nodeid.NodeID,
    signingPrivkey *fcrcrypto.KeyPair,
    signingPrivKeyVer *fcrcrypto.KeyVersion,
  ) error

  RequestSetReputation(
    gatewayRegistrar register.GatewayRegistrar,
    clientID *nodeid.NodeID,
    reputation int64,
    signingPrivkey *fcrcrypto.KeyPair,
    signingPrivKeyVer *fcrcrypto.KeyVersion,
  ) (bool, error)
}

func NewAdminApi() AdminApi {
  return &Admin{
    httpCommunicator: request.NewHttpCommunicator(),
  }
}

func NewAdminApiWithDep(httpCommunicator request.HttpCommunications) AdminApi {
  return &Admin{
    httpCommunicator: httpCommunicator,
  }
}
