package adminapi

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
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
		providerRegistrar register.ProviderRegistrar,
		signingPrivkey *fcrcrypto.KeyPair,
		signingPrivKeyVer *fcrcrypto.KeyVersion,
	) error

	RequestGetPublishedOffer(
		providerRegistrar register.ProviderRegistrar,
		gatewayIDs []nodeid.NodeID,
		signingPrivkey *fcrcrypto.KeyPair,
		signingPrivKeyVer *fcrcrypto.KeyVersion,
	) (
		bool,                // found
		[]cidoffer.CIDOffer, // offers
		error,               // error
	)

	RequestInitialiseKey(
		providerRegistrar register.ProviderRegistrar,
		providerPrivKey *fcrcrypto.KeyPair,
		providerPrivKeyVer *fcrcrypto.KeyVersion,
		signingPrivkey *fcrcrypto.KeyPair,
		signingPrivKeyVer *fcrcrypto.KeyVersion,
	) error

	RequestInitialiseKeyV2(
		providerRegistrar register.ProviderRegistrar,
		providerPrivKey *fcrcrypto.KeyPair,
		providerPrivKeyVer *fcrcrypto.KeyVersion,
		signingPrivkey *fcrcrypto.KeyPair,
		signingPrivKeyVer *fcrcrypto.KeyVersion,
		lotusWalletPrivateKey string,
		lotusAP string,
		lotusAuthToken string,
	) error

	RequestPublishDHTOffer(
		providerRegistrar register.ProviderRegistrar,
		cids []cid.ContentID,
		price []uint64,
		expiry []int64,
		qos []uint64,
		signingPrivkey *fcrcrypto.KeyPair,
		signingPrivKeyVer *fcrcrypto.KeyVersion,
	) error

	RequestPublishGroupOffer(
		providerRegistrar register.ProviderRegistrar,
		cids []cid.ContentID,
		price uint64,
		expiry int64,
		qos uint64,
		signingPrivkey *fcrcrypto.KeyPair,
		signingPrivKeyVer *fcrcrypto.KeyVersion,
	) error
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