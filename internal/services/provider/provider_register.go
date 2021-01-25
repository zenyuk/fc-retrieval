package provider

import (
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-provider/internal/request"
)

func register(url string) {

	providerReg := Register{
		Address:        "f01213",
		NetworkInfo:    "127.0.0.1:80",
		RegionCode:     "US",
		RootSigningKey: "0xABCDE123456789",
		SigingKey:      "0x987654321EDCBA",
	}

	err := request.PostJSON(url, providerReg)
	if err != nil {
		log.Error("%+v", err)
	}
}
