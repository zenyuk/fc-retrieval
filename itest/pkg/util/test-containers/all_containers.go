package test_containers

type AllContainers struct {
	Gateways  map[string]GatewayPortsResolver
	Providers map[string]ProviderPortsResolver
	Register  RegisterPortsResolver
	Redis     Container
	Lotus     LotusPortsResolver
}
