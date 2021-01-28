package provider

import (
	"time"

	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-provider/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-provider/internal/register"
	"github.com/ConsenSys/fc-retrieval-provider/pkg/communication"
	"github.com/spf13/viper"
)

// Provider configuration
type Provider struct {
	Conf 						*viper.Viper
	GatewayCommPool *communication.CommunicationPool
}

// NewProvider returns new provider
func NewProvider(conf *viper.Viper) *Provider {
	gatewayCommPool := communication.NewCommunicationPool()
	return &Provider{
		Conf: 						conf,
		GatewayCommPool: 	&gatewayCommPool,
	}
}

// Start a provider
func (provider *Provider) Start() {
	provider.greet()
	provider.registration()
	provider.loop()
}

// Greeting
func (provider *Provider) greet() {
	scheme := provider.Conf.GetString("SERVICE_SCHEME")
	host := provider.Conf.GetString("SERVICE_HOST")
	port := provider.Conf.GetString("SERVICE_PORT")
	log.Info("Provider started at %s://%s:%s", scheme, host, port)
}

// Register the provider
func (provider *Provider) registration() {
	err := register.RegisterProvider(provider.Conf)
	if err != nil {
		log.Error("Provider not registered: %v", err)
		//TODO graceful exit
	}
}

// Start infinite loop
func (provider *Provider) loop() {
	for {
		gateways, err := register.GetRegisteredGateways(provider.Conf)
		if err != nil {
			log.Error("Unable to get registered gateways: %v", err)
			//TODO graceful exit
		}
		for _, gw := range gateways {
			message := generateDummyMessage()
			log.Info("Message: %v", message)
			
			// TODO: remove
			gw.NodeID = "101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2D2E2F"
			// End TODO

			gatewayID, err := nodeid.NewNodeIDFromString(gw.NodeID)
			if err != nil {
				log.Error("Error with nodeID %v: %v", gw.NodeID, err)
				continue
			}
			provider.GatewayCommPool.RegisterNodeAddress(gatewayID, gw.NetworkInfo)
			gateway.SendMessage(message, gatewayID, provider.GatewayCommPool)
		}
		time.Sleep(25 * time.Second)
	}
}
