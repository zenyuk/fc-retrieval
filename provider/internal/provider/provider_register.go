package provider

import (
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-provider/internal/request"
)

func register(url string, p *Provider) {

	providerReg := Register{
		Address:        p.conf.GetString("PROVIDER_ADDRESS"),
		NetworkInfo:    p.conf.GetString("PROVIDER_NETWORK_INFO"),
		RegionCode:     p.conf.GetString("PROVIDER_REGION_CODE"),
		RootSigningKey: p.conf.GetString("PROVIDER_ROOT_SIGNING_KEY"),
		SigingKey:      p.conf.GetString("PROVIDER_SIGNING_KEY"),
	}

	err := request.SendJSON(url, providerReg)
	if err != nil {
		log.Error("%+v", err)
	}
}
