package provider

import (
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

	request.PostJSON(url, providerReg)
}
