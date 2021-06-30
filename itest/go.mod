module github.com/ConsenSys/fc-retrieval/itest

go 1.16

replace github.com/ConsenSys/fc-retrieval/common => ../common

replace github.com/ConsenSys/fc-retrieval/client => ../client

replace github.com/ConsenSys/fc-retrieval/gateway-admin => ../gateway-admin

replace github.com/ConsenSys/fc-retrieval/provider-admin => ../provider-admin

require (
	github.com/ConsenSys/fc-retrieval/client v0.0.0-00010101000000-000000000000
	github.com/ConsenSys/fc-retrieval/common v0.0.0-00010101000000-000000000000
	github.com/ConsenSys/fc-retrieval/gateway-admin v0.0.0-00010101000000-000000000000
	github.com/ConsenSys/fc-retrieval/provider-admin v0.0.0-00010101000000-000000000000
	github.com/c-bata/go-prompt v0.2.6
	github.com/davidlazar/go-crypto v0.0.0-20200604182044-b73af7476f6c // indirect
	github.com/docker/docker v20.10.5+incompatible
	github.com/filecoin-project/go-address v0.0.5
	github.com/filecoin-project/go-jsonrpc v0.1.4-0.20210217175800-45ea43ac2bec
	github.com/filecoin-project/lotus v1.8.0
	github.com/filecoin-project/specs-actors/v4 v4.0.1 // indirect
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/google/uuid v1.2.0
	github.com/ipfs/go-cid v0.0.7
	github.com/kr/text v0.2.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/pelletier/go-toml v1.7.0 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/testcontainers/testcontainers-go v0.10.0
	github.com/wcgcyx/testcontainers-go v0.10.1-0.20210511154849-504eecefabe0
	github.com/whyrusleeping/cbor-gen v0.0.0-20210303213153-67a261a1d291 // indirect
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
