module github.com/ConsenSys/fc-retrieval-itest

go 1.15

require (
	github.com/ConsenSys/fc-retrieval-client v0.0.0-20210315220115-b5a58266695e
	github.com/ConsenSys/fc-retrieval-common v0.0.0-20210312151557-4caead038a43
	github.com/ConsenSys/fc-retrieval-gateway-admin v0.0.0-20210317222244-2058ef0b3d96
	github.com/ConsenSys/fc-retrieval-provider-admin v0.0.0-20210315220609-1fe0ee54f441
	github.com/ConsenSys/fc-retrieval-register v0.0.0-20210315215728-57ff758e2e2c
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
)

replace github.com/ConsenSys/fc-retrieval-provider-admin => ../fc-retrieval-provider-admin
replace github.com/ConsenSys/fc-retrieval-gateway-admin => ../fc-retrieval-gateway-admin
replace github.com/ConsenSys/fc-retrieval-common => ../fc-retrieval-common
replace github.com/ConsenSys/fc-retrieval-client => ../fc-retrieval-client