module github.com/ConsenSys/fc-retrieval-itest

go 1.15

require (
	github.com/ConsenSys/fc-retrieval-client v0.0.0-20210331013001-37693cd8732f
	github.com/ConsenSys/fc-retrieval-common v0.0.0-20210331012500-bf32c4270d66
	github.com/ConsenSys/fc-retrieval-gateway-admin v0.0.0-20210331013355-cce2ba8c3ff2
	github.com/ConsenSys/fc-retrieval-provider-admin v0.0.0-20210331013817-cfbd60705f9e
	github.com/ConsenSys/fc-retrieval-register v0.0.0-20210331012722-2659896d78fb
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
)

replace github.com/ConsenSys/fc-retrieval-provider-admin => ../fc-retrieval-provider-admin
replace github.com/ConsenSys/fc-retrieval-gateway-admin => ../fc-retrieval-gateway-admin
replace github.com/ConsenSys/fc-retrieval-common => ../fc-retrieval-common
replace github.com/ConsenSys/fc-retrieval-client => ../fc-retrieval-client