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
	conf 			*viper.Viper
	gCommPool *communication.CommunicationPool
}

// NewProvider returns new provider
func NewProvider(conf *viper.Viper) *Provider {
	gCommPool := communication.NewCommunicationPool()
	return &Provider{
		conf: 			conf,
		gCommPool: 	&gCommPool,
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
	scheme := provider.conf.GetString("SERVICE_SCHEME")
	host := provider.conf.GetString("SERVICE_HOST")
	port := provider.conf.GetString("SERVICE_PORT")
	log.Info("Provider started at %s://%s:%s", scheme, host, port)
}

// Register the provider
func (provider *Provider) registration() {
	err := register.RegisterProvider(provider.conf)
	if err != nil {
		log.Error("Provider not registered: %v", err)
		//TODO exit
	}
}

// Start infinite loop
func (provider *Provider) loop() {
	for {
		gateways, err := register.GetRegisteredGateways(provider.conf)
		if err != nil {
			log.Error("Unable to get registered gateways: %v", err)
			//TODO exit
		}
		for _, gw := range gateways {
			message := generateDummyMessage()
			log.Info("Message: %v", message)
			gCommPool := provider.gCommPool
			// TODO: remove
			gw.NodeID = "101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2D2E2F"
			gatewayID, _ := nodeid.NewNodeIDFromString(gw.NodeID)
			// End TODO
			gCommPool.RegisterNodeAddress(gatewayID, gw.NetworkInfo)
			gateway.SendMessage(message, gatewayID, gCommPool)
		}
		time.Sleep(25 * time.Second)
	}
}
