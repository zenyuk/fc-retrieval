package provider

import (
	"fmt"
	"time"

	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
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
	p.loop()
}

func (p *Provider) loop() {
	url := p.conf.GetString("REGISTER_API_URL") + "/registers/gateway"

	register(url)
	for {
		gateways := []Register{}
		request.GetJSON(url, &gateways)

		for _, gateway := range gateways {
			message := generateDummyMessage()
			fmt.Println("TODO, post to this gateway")
			fmt.Println(gateway)
			fmt.Println(message)
		}
		time.Sleep(25 * time.Second)
	}
}
