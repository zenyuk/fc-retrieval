package provider

import (
	"math/rand"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

func generateDummyMessage() messages.ProviderPublishGroupCIDRequest {
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	protocolSupported := []int32{1, 2}
	providerID, _ := nodeid.NewRandomNodeID()

	contentID, _ := cid.NewRandomContentID()
	pieceCIDs := []cid.ContentID{*contentID}
	dummyMessage := messages.ProviderPublishGroupCIDRequest{

		MessageType:       123,
		ProtocolVersion:   1,
		ProtocolSupported: protocolSupported,
		Nonce:             rand.Int63n(100000),
		ProviderID:        *providerID,
		Price:             42,
		Expiry:            expiryDate,
		QoS:               42,
		Signature:         "Signature",
		PieceCIDs:         pieceCIDs,
	}
	return dummyMessage
}
