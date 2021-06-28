module github.com/ConsenSys/fc-retrieval/itest

go 1.16

replace github.com/ConsenSys/fc-retrieval/common => ../common

require (
	github.com/ConsenSys/fc-retrieval/client v0.0.0-00010101000000-000000000000
	github.com/ConsenSys/fc-retrieval/common v0.0.0-00010101000000-000000000000
	github.com/ConsenSys/fc-retrieval/gateway-admin v0.0.0-00010101000000-000000000000
	github.com/ConsenSys/fc-retrieval/provider-admin v0.0.0-00010101000000-000000000000
	github.com/docker/docker v20.10.5+incompatible
	github.com/filecoin-project/go-address v0.0.5
	github.com/filecoin-project/go-jsonrpc v0.1.4-0.20210217175800-45ea43ac2bec
	github.com/filecoin-project/lotus v1.8.0
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/uuid v1.2.0
	github.com/ipfs/go-cid v0.0.7
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/testcontainers/testcontainers-go v0.10.0
	github.com/wcgcyx/testcontainers-go v0.10.1-0.20210511154849-504eecefabe0
)

replace github.com/ConsenSys/fc-retrieval/client => ../client

replace github.com/ConsenSys/fc-retrieval/gateway-admin => ../gateway-admin

replace github.com/ConsenSys/fc-retrieval/provider-admin => ../provider-admin
