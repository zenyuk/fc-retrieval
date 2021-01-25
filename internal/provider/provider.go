package provider

import (
	"fmt"
	"time"

	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-provider/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-provider/internal/request"
	"github.com/spf13/viper"
)

// Provider configuration
type Provider struct {
	conf *viper.Viper
}

// Register data model
type Register struct {
	Address        string
	NetworkInfo    string
	RegionCode     string
	RootSigningKey string
	SigingKey      string
}

// NewProvider returns new provider
func NewProvider(conf *viper.Viper) *Provider {
	return &Provider{
		conf: conf,
	}
}

// Start a new provider
func Start(p *Provider) {
	scheme := p.conf.GetString("SERVICE_SCHEME")
	host := p.conf.GetString("SERVICE_HOST")
	port := p.conf.GetString("SERVICE_PORT")
	log.Info("Provider started at %s://%s:%s", scheme, host, port)
	url := p.conf.GetString("REGISTER_API_URL") + "/registers/provider"
	register(url)

	// Temp register gateway
	urlGw := p.conf.GetString("REGISTER_API_URL") + "/registers/gateway"
	register(urlGw)
	p.loop()
}

// Start infinite loop
func (p *Provider) loop() {
	url := p.conf.GetString("REGISTER_API_URL") + "/registers/gateway"

	// connect to gateway and store connection
	for {
		log.Info(".")
		gateways := []Register{}
		request.GetJSON(url, &gateways)

		for _, gw := range gateways {
			message := generateDummyMessage()
			fmt.Println(message)
			gwURL := gw.NetworkInfo
			gateway.SendMessage(gwURL, message)
		}
		time.Sleep(25 * time.Second)
	}
}
