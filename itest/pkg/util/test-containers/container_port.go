package test_containers

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/spf13/viper"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
)

type ContainerPorts struct {
	ctx              context.Context
	container        tc.Container
	guestToHostPorts map[int]nat.Port
	gatewayConfig    *viper.Viper
	providerConfig   *viper.Viper
	registerConfig   *viper.Viper
}

type Container interface {
	TerminateContainer()
	SerialisePortSetup() string
}

type RegisterPortsResolver interface {
	Container
	GetRegisterHostApiEndpoint() (registerApiEndpoint string)
}

type GatewayPortsResolver interface {
	Container
	GetGatewayHostApiEndpoints() (gatewayApiEndpoint, providerApiEndpoint, clientApiEndpoint, adminApiEndpoint string)
}

type ProviderPortsResolver interface {
	Container
	GetProviderHostApiEndpoints() (gatewayApiEndpoint, clientApiEndpoint, adminApiEndpoint string)
}

type LotusPortsResolver interface {
	Container
	GetLostHostApiEndpoints() (lotusDaemonApiEndpoint string, lotusMinerApiEndpoint string)
}

func NewGenericPortsResolver(ctx context.Context, container tc.Container, guestToHostPorts map[int]nat.Port) Container {
	return ContainerPorts{
		ctx,
		container,
		guestToHostPorts,
		nil,
		nil,
		nil,
	}
}

func NewGatewayPortsResolver(ctx context.Context, container tc.Container, guestToHostPorts map[int]nat.Port, config *viper.Viper) GatewayPortsResolver {
	return ContainerPorts{
		ctx,
		container,
		guestToHostPorts,
		config,
		nil,
		nil,
	}
}

func NewProviderPortsResolver(ctx context.Context, container tc.Container, guestToHostPorts map[int]nat.Port, config *viper.Viper) ProviderPortsResolver {
	return ContainerPorts{
		ctx,
		container,
		guestToHostPorts,
		nil,
		config,
		nil,
	}
}

func NewRegisterPortsResolver(ctx context.Context, container tc.Container, guestToHostPorts map[int]nat.Port, config *viper.Viper) RegisterPortsResolver {
	return ContainerPorts{
		ctx,
		container,
		guestToHostPorts,
		nil,
		nil,
		config,
	}
}

func NewLotusPortsResolver(ctx context.Context, container tc.Container, guestToHostPorts map[int]nat.Port) LotusPortsResolver {
	return ContainerPorts{
		ctx,
		container,
		guestToHostPorts,
		nil,
		nil,
		nil,
	}
}

func (c ContainerPorts) SerialisePortSetup() string {
	containerName, err := c.container.Name(c.ctx)
	if err != nil {
		logging.Fatal("GetGatewayHostApiEndpoints - container name resolution error: %s", err)
	}
	result := fmt.Sprintf(">> container: %-11s ports: ", containerName)
	for guestPort, hostPort := range c.guestToHostPorts {
		result += fmt.Sprintf("%s->%d ", hostPort.Port(), guestPort)
	}
	result += "\n"
	return result
}

func (c ContainerPorts) GetRegisterHostApiEndpoint() (registerApiEndpoint string) {
	// expecting containerName with extra slash char in the beginning, like: "/bla"
	containerName, err := c.container.Name(c.ctx)
	if err != nil {
		logging.Fatal("GetGatewayHostApiEndpoints - container name resolution error: %s", err)
	}
	registerApiEndpoint = containerName[1:] + ":" + c.guestToHostPorts[c.registerConfig.GetInt("BIND_API")].Port()
	return
}

func (c ContainerPorts) GetGatewayHostApiEndpoints() (gatewayApiEndpoint, providerApiEndpoint, clientApiEndpoint, adminApiEndpoint string) {
	containerName, err := c.container.Name(c.ctx)
	if err != nil {
		logging.Fatal("GetGatewayHostApiEndpoints - container name resolution error: %s", err)
	}
	// expecting containerName with extra slash char in the beginning, like: "/bla"
	containerName = containerName[1:]
	gatewayApiEndpoint = containerName + ":" + c.guestToHostPorts[c.gatewayConfig.GetInt("BIND_GATEWAY_API")].Port()
	providerApiEndpoint = containerName + ":" + c.guestToHostPorts[c.gatewayConfig.GetInt("BIND_PROVIDER_API")].Port()
	clientApiEndpoint = containerName + ":" + c.guestToHostPorts[c.gatewayConfig.GetInt("BIND_REST_API")].Port()
	adminApiEndpoint = containerName + ":" + c.guestToHostPorts[c.gatewayConfig.GetInt("BIND_ADMIN_API")].Port()
	return
}

func (c ContainerPorts) GetProviderHostApiEndpoints() (gatewayApiEndpoint, clientApiEndpoint, adminApiEndpoint string) {
	containerName, err := c.container.Name(c.ctx)
	if err != nil {
		logging.Fatal("GetProviderHostApiEndpoints - container name resolution error: %s", err)
	}
	containerName = containerName[1:]
	gatewayApiEndpoint = containerName + ":" + c.guestToHostPorts[c.providerConfig.GetInt("BIND_GATEWAY_API")].Port()
	clientApiEndpoint = containerName + ":" + c.guestToHostPorts[c.providerConfig.GetInt("BIND_REST_API")].Port()
	adminApiEndpoint = containerName + ":" + c.guestToHostPorts[c.providerConfig.GetInt("BIND_ADMIN_API")].Port()
	return
}

func (c ContainerPorts) GetLostHostApiEndpoints() (lotusDaemonApiEndpoint string, lotusMinerApiEndpoint string) {
	containerName, err := c.container.Name(c.ctx)
	if err != nil {
		logging.Fatal("GetProviderHostApiEndpoints - container name resolution error: %s", err)
	}
	containerName = containerName[1:]
	lotusDaemonApiEndpoint = containerName + ":" + c.guestToHostPorts[1234].Port()
	lotusMinerApiEndpoint = containerName + ":" + c.guestToHostPorts[2345].Port()
	return
}

func (c ContainerPorts) TerminateContainer() {
	if c.container == nil {
		// already terminated, do nothing
		return
	}

	if err := c.container.StopLogProducer(); err != nil {
		logging.Error("error stopping logs for a test container: %s", err.Error())
	}
	if err := c.container.Terminate(c.ctx); err != nil {
		logging.Fatal("error terminating test container: %s", err.Error())
	}
}
