package provider

import (
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-provider/internal/request"
)

func register(url string) {

	providerReg := Register{
		Address:        "f0121345",
		NetworkInfo:    "127.0.0.1:8090",
		RegionCode:     "US",
		RootSigningKey: "0xABCDE123456789",
		SigingKey:      "0x987654321EDCBA",
	}

	err := request.SendJSON(url, providerReg)
	if err != nil {
		log.Error("%+v", err)
	}
}
