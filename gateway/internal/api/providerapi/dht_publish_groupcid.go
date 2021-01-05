package providerapi

import (
	"encoding/json"
	"net"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

func handleProviderDHTPublishGroupCIDRequest(conn net.Conn, request *messages.ProviderDHTPublishGroupCIDRequest) error {
	// Do something about the internal request
	// Will need gateway instance to read from db
	// gateway, err := gateway.GetSingleInstance()
	logging.Info("Provider request from: %s", request.ProviderID.ToString())
	// Respond to provider.
	response, _ := json.Marshal(messages.ProviderDHTPublishGroupCIDResponse{
		MessageType: messages.GatewayDHTDiscoverResponseType,
		// This is just a dummy response
	})
	return tcpcomms.SendTCPMessage(conn, messages.ProviderDHTPublishGroupCIDResponseType, response, settings.DefaultTCPInactivityTimeoutMs*time.Millisecond)
}
