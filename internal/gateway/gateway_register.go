package gateway

import (
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/request"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
)

// Registration for Gateway
func Registration(url string, settings settings.AppSettings) {

	providerReg := register.GatewayRegister{
		NodeID:              settings.GatewayID,
		Address:             settings.GatewayAddress,
		NetworkGatewayInfo:  settings.GatewayNetworkInfo,
		NetworkProviderInfo: settings.GatewayNetworkInfo,
		NetworkClientInfo:   settings.GatewayNetworkInfo,
		NetworkAdminInfo:    settings.GatewayNetworkInfo,
		RegionCode:          settings.GatewayRegionCode,
		RootSigningKey:      settings.GatewayRootSigningKey,
		SigingKey:           settings.GatewaySigningKey,
	}

	err := request.SendJSON(url, providerReg)
	if err != nil {
		log.Error("%+v", err)
	}
}
