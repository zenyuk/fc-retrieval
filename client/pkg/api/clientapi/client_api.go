package clientapi

import (
	"github.com/ConsenSys/fc-retrieval/common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval/common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval/common/pkg/register"
	"github.com/ConsenSys/fc-retrieval/common/pkg/request"
)

type Client struct {
	httpCommunicator request.HttpCommunications
}

type ClientApi interface {
	RequestDHTOfferDiscover(
		clientApiEndpoint string,
		gatewayInfo register.GatewayRegistrar,
		gatewayIDs []nodeid.NodeID,
		contentID *cid.ContentID,
		nonce int64,
		offersDigests [][][cidoffer.CIDOfferDigestSize]byte,
		paymentChannelAddr string,
		voucher string,
	) ([]GatewaySubOffers, error)

	RequestDHTDiscover(
		clientApiEndpoint string,
		gatewayInfo register.GatewayRegistrar,
		contentID *cid.ContentID,
		nonce int64,
		ttl int64,
		numDHT int64,
		incrementalResult bool,
		paychAddr string,
		voucher string,
	) ([]nodeid.NodeID, []fcrmessages.FCRMessage, []nodeid.NodeID, error)

	RequestDHTDiscoverV2(
		clientApiEndpoint string,
		gatewayInfo register.GatewayRegistrar,
		contentID *cid.ContentID,
		nonce int64,
		ttl int64,
		numDHT int64,
		incrementalResult bool,
		paychAddr string,
		voucher string,
	) ([]nodeid.NodeID, []fcrmessages.FCRMessage, []nodeid.NodeID, bool, string, error)

	RequestDHTOfferAck(
		clientApiEndpoint string,
		providerInfo register.ProviderRegistrar,
		contentID *cid.ContentID,
		gatewayID *nodeid.NodeID,
	) (bool, *fcrmessages.FCRMessage, *fcrmessages.FCRMessage, error)

	RequestEstablishment(
		clientApiEndpoint string,
		gatewayInfo register.GatewayRegistrar,
		challenge []byte,
		clientID *nodeid.NodeID,
		ttl int64,
	) error

	RequestStandardDiscoverOffer(
		clientApiEndpoint string,
		gatewayInfo register.GatewayRegistrar,
		contentID *cid.ContentID,
		nonce int64,
		ttl int64,
		offerDigests [][cidoffer.CIDOfferDigestSize]byte,
		paychAddr string,
		voucher string,
	) ([]cidoffer.SubCIDOffer, error)

	RequestStandardDiscover(
		clientApiEndpoint string,
		gatewayInfo register.GatewayRegistrar,
		contentID *cid.ContentID,
		nonce int64,
		ttl int64,
		paychAddr string,
		voucher string,
	) ([]cidoffer.SubCIDOffer, error)

	RequestStandardDiscoverV2(
		clientApiEndpoint string,
		gatewayInfo register.GatewayRegistrar,
		contentID *cid.ContentID,
		nonce int64,
		ttl int64,
		paychAddr string,
		voucher string,
	) ([][cidoffer.CIDOfferDigestSize]byte, bool, string, error)
}

func NewClientApi() ClientApi {
	return &Client{
		httpCommunicator: request.NewHttpCommunicator(),
	}
}

func NewAdminApiWithDep(httpCommunicator request.HttpCommunications) ClientApi {
	return &Client{
		httpCommunicator: httpCommunicator,
	}
}
