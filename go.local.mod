module github.com/ConsenSys/fc-retrieval-itest
go 1.15
require (
	github.com/ConsenSys/fc-retrieval-client v0.0.0-20210310062823-1492c15b2fab
	github.com/ConsenSys/fc-retrieval-common v0.0.0-20210309021945-823304bbc3fc
	github.com/ConsenSys/fc-retrieval-gateway-admin v0.0.0-20210309023153-45f2599e0598
	github.com/ConsenSys/fc-retrieval-provider-admin v0.0.0-20210309023212-0c2328a0d2a1
	github.com/ConsenSys/fc-retrieval-register v0.0.0-20210310004700-1507a8b3268f
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
)
replace github.com/ConsenSys/fc-retrieval-provider-admin => ../fc-retrieval-provider-admin
replace github.com/ConsenSys/fc-retrieval-gateway-admin => ../fc-retrieval-gateway-admin
replace github.com/ConsenSys/fc-retrieval-common => ../fc-retrieval-common
replace github.com/ConsenSys/fc-retrieval-client => ../fc-retrieval-client