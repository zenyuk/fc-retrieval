package test_containers

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/ConsenSys/fc-retrieval/common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval/itest/pkg/util/test-containers/color"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
)

const lotusFullNodeWaitFor = "starting winning PoSt warmup"
const networkMode = "default"

func StartContainers(ctx context.Context, gatewaysCount int, providersCount int, testCaseName string, verbose bool, gatewayConfig *viper.Viper, providerConfig *viper.Viper, registerConfig *viper.Viper) (
	containers AllContainers,
	network *tc.Network,
	err error) {

	if gatewaysCount < 1 || providersCount < 1 {
		return containers, nil, errors.New("number of gateways or providers in a test setup can't be less then one")
	}
	network, networkName := createNetwork(ctx)
	containers.Redis, err = startRedis(ctx, networkName, verbose)
	if err != nil {
		return containers, nil, fmt.Errorf("can't start Redis container, error: %s", err.Error())
	}
	containers.Lotus, err = startLotusFullNode(ctx, networkName, verbose)
	if err != nil {
		return containers, nil, fmt.Errorf("can't start Lotus container, error: %s", err.Error())
	}
	containers.Register, err = startRegister(ctx, networkName, verbose, registerConfig)
	if err != nil {
		return containers, nil, fmt.Errorf("can't start Register container, error: %s", err.Error())
	}

	// Start all providers
	containers.Providers = make(map[string]ProviderPortsResolver)
	for i := 0; i < providersCount; i++ {
		providerName := fmt.Sprintf("provider-%v", i)
		containers.Providers[providerName], err = startProvider(ctx, providerName, networkName, verbose, providerConfig)
		if err != nil {
			return containers, nil, fmt.Errorf("can't start Provider container, error: %s", err.Error())
		}
	}
	// Start all gateways
	containers.Gateways = make(map[string]GatewayPortsResolver)
	for i := 0; i < gatewaysCount; i++ {
		gatewayName := fmt.Sprintf("gateway-%v", i)
		containers.Gateways[gatewayName], err = startGateway(ctx, gatewayName, networkName, verbose, gatewayConfig)
		if err != nil {
			return containers, nil, fmt.Errorf("can't start Gateway container, error: %s", err.Error())
		}
	}
	printPortsSetup(testCaseName, networkName, containers)
	return
}

func StopContainers(ctx context.Context, testPackage string, containers AllContainers, network *tc.Network) {
	fmt.Printf("\n>> terminating containers for test package: %s\n\n", testPackage)
	for _, g := range containers.Gateways {
		g.TerminateContainer()
	}
	for _, p := range containers.Providers {
		p.TerminateContainer()
	}
	containers.Register.TerminateContainer()
	containers.Lotus.TerminateContainer()
	containers.Redis.TerminateContainer()
	if err := (*network).Remove(ctx); err != nil {
		logging.Error("error terminating test container network: %s", err.Error())
	}
}

// getEnvMap gets the env map from a given env file
func getEnvMap(envFile string) map[string]string {
	env := make(map[string]string)
	file, err := os.Open(envFile)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res := strings.Split(scanner.Text(), "=")
		if len(res) >= 2 {
			env[res[0]] = res[1]
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if err := file.Close(); err != nil {
		logging.Error("error closing file: %s", err.Error())
	}
	return env
}

// CreateNetwork creates a network
func createNetwork(ctx context.Context) (*tc.Network, string) {
	randomUuid, _ := uuid.NewRandom()
	networkName := randomUuid.String()
	net, err := tc.GenericNetwork(ctx, tc.GenericNetworkRequest{
		NetworkRequest: tc.NetworkRequest{
			Name:           networkName,
			CheckDuplicate: true,
		},
	})
	if err != nil {
		panic(err)
	}
	return &net, networkName
}

// startRedis - starts redis; used by register
func startRedis(ctx context.Context, network string, verbose bool) (Container, error) {
	guestPort, _ := nat.NewPort("tcp", "6379")
	containerName := "redis"
	logConsumer := TestLogConsumer{
		Messages:      []string{},
		ContainerName: containerName,
		Color:         color.White,
	}
	req := tc.ContainerRequest{
		Name:         containerName,
		Image:        "redis:alpine",
		Cmd:          []string{"redis-server", "--requirepass", "xxxx"},
		ExposedPorts: []string{guestPort.Port()},
		//Env:            map[string]string{"ALLOW_EMPTY_PASSWORD": "yes"},
		WaitingFor:     wait.ForLog("Ready to accept connections"),
		AutoRemove:     true,
		Networks:       []string{network},
		NetworkMode:    container.NetworkMode(networkMode),
		NetworkAliases: map[string][]string{network: {"redis"}},
	}
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}

	if verbose {
		container.StartLogProducer(ctx)
		container.FollowOutput(&logConsumer)
	}

	hostPort, err := getMappedHostPort(ctx, guestPort, container)
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}
	return NewGenericPortsResolver(ctx, container, map[int]nat.Port{guestPort.Int(): hostPort}), nil
}

func startLotusFullNode(ctx context.Context, network string, verbose bool) (LotusPortsResolver, error) {
	guestPort1, _ := nat.NewPort("tcp", "1234")
	guestPort2, _ := nat.NewPort("tcp", "2345")
	containerName := "lotus"
	logConsumer := TestLogConsumer{
		Messages:      []string{},
		ContainerName: containerName,
		Color:         color.Gray,
	}
	req := tc.ContainerRequest{
		Name:           containerName,
		Image:          "consensys/lotus-full-node:latest",
		ExposedPorts:   []string{guestPort1.Port(), guestPort2.Port()},
		Networks:       []string{network},
		NetworkMode:    container.NetworkMode(networkMode),
		NetworkAliases: map[string][]string{network: {containerName}},
		WaitingFor:     wait.ForLog(lotusFullNodeWaitFor).WithStartupTimeout(10 * time.Minute),
		AutoRemove:     true,
	}
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}

	if verbose {
		container.StartLogProducer(ctx)
		container.FollowOutput(&logConsumer)
	}

	hostPort1, err := getMappedHostPort(ctx, guestPort1, container)
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}
	hostPort2, err := getMappedHostPort(ctx, guestPort2, container)
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}
	return NewLotusPortsResolver(ctx, container, map[int]nat.Port{guestPort1.Int(): hostPort1, guestPort2.Int(): hostPort2}), nil
}

func startRegister(ctx context.Context, network string, verbose bool, config *viper.Viper) (RegisterPortsResolver, error) {
	env := getEnvMap(config.ConfigFileUsed())
	guestPort, _ := nat.NewPort("tcp", "9020")
	containerName := "register"
	logConsumer := TestLogConsumer{
		Messages:      []string{},
		ContainerName: containerName,
		Color:         color.Yellow,
	}
	req := tc.ContainerRequest{
		Name:           containerName,
		Image:          "consensys/fc-retrieval/register:latest",
		ExposedPorts:   []string{guestPort.Port()},
		Networks:       []string{network},
		Env:            env,
		NetworkMode:    container.NetworkMode(networkMode),
		NetworkAliases: map[string][]string{network: {containerName}},
		WaitingFor:     wait.ForLog("Serving register at"),
		AutoRemove:     true,
	}
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}

	if verbose {
		container.StartLogProducer(ctx)
		container.FollowOutput(&logConsumer)
	}

	hostPort, err := getMappedHostPort(ctx, guestPort, container)
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}
	return NewRegisterPortsResolver(ctx, container, map[int]nat.Port{guestPort.Int(): hostPort}, config), nil
}

func startGateway(ctx context.Context, gatewayName string, network string, verbose bool, config *viper.Viper) (GatewayPortsResolver, error) {
	env := getEnvMap(config.ConfigFileUsed())
	guestPort1, _ := nat.NewPort("tcp", "9010")
	guestPort2, _ := nat.NewPort("tcp", "9011")
	guestPort3, _ := nat.NewPort("tcp", "9012")
	guestPort4, _ := nat.NewPort("tcp", "9013")
	logConsumer := TestLogConsumer{
		Messages:      []string{},
		ContainerName: gatewayName,
		Color:         color.Green,
	}
	req := tc.ContainerRequest{
		Name:           gatewayName,
		Image:          "consensys/fc-retrieval/gateway:latest",
		ExposedPorts:   []string{guestPort1.Port(), guestPort2.Port(), guestPort3.Port(), guestPort4.Port()},
		Networks:       []string{network},
		Env:            env,
		NetworkMode:    container.NetworkMode(networkMode),
		NetworkAliases: map[string][]string{network: {gatewayName}},
		WaitingFor:     wait.ForLog("Filecoin Gateway Start-up Complete"),
		AutoRemove:     true,
	}
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}

	if verbose {
		container.StartLogProducer(ctx)
		container.FollowOutput(&logConsumer)
	}

	hostPort1, err := getMappedHostPort(ctx, guestPort1, container)
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}
	hostPort2, err := getMappedHostPort(ctx, guestPort2, container)
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}
	hostPort3, err := getMappedHostPort(ctx, guestPort3, container)
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}
	hostPort4, err := getMappedHostPort(ctx, guestPort4, container)
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}
	return NewGatewayPortsResolver(ctx, container, map[int]nat.Port{guestPort1.Int(): hostPort1, guestPort2.Int(): hostPort2, guestPort3.Int(): hostPort3, guestPort4.Int(): hostPort4}, config), nil
}

func startProvider(ctx context.Context, providerName string, network string, verbose bool, config *viper.Viper) (ProviderPortsResolver, error) {
	env := getEnvMap(config.ConfigFileUsed())
	guestPort1, _ := nat.NewPort("tcp", "9030")
	guestPort2, _ := nat.NewPort("tcp", "9032")
	guestPort3, _ := nat.NewPort("tcp", "9033")
	logConsumer := TestLogConsumer{
		Messages:      []string{},
		ContainerName: providerName,
		Color:         color.Purple,
	}

	req := tc.ContainerRequest{
		Name:           providerName,
		Image:          "consensys/fc-retrieval/provider:latest",
		ExposedPorts:   []string{guestPort1.Port(), guestPort2.Port(), guestPort3.Port()},
		Networks:       []string{network},
		Env:            env,
		NetworkMode:    container.NetworkMode(networkMode),
		NetworkAliases: map[string][]string{network: {providerName}},
		WaitingFor:     wait.ForLog("Filecoin Provider Start-up Complete"),
		AutoRemove:     true,
	}
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}

	if verbose {
		container.StartLogProducer(ctx)
		container.FollowOutput(&logConsumer)
	}

	hostPort1, err := getMappedHostPort(ctx, guestPort1, container)
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}
	hostPort2, err := getMappedHostPort(ctx, guestPort2, container)
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}
	hostPort3, err := getMappedHostPort(ctx, guestPort3, container)
	if err != nil {
		if container != nil {
			container.Terminate(ctx)
		}
		return nil, err
	}
	return NewProviderPortsResolver(ctx, container, map[int]nat.Port{guestPort1.Int(): hostPort1, guestPort2.Int(): hostPort2, guestPort3.Int(): hostPort3}, config), nil
}

// getMappedHostPort returns the actual host port. The result port is dynamically allocated by TestContainers
func getMappedHostPort(ctx context.Context, guestPort nat.Port, container tc.Container) (nat.Port, error) {
	logging.Debug("getMappedHostPort - context is alive: %+v; guestPort: %s; container ID: %s\n", ctx != nil, guestPort.Port(), container.GetContainerID()[:12])
	_, err := container.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("can't resolve host for the given container ID: %s; error: %s", container.GetContainerID()[:12], err.Error())
	}

	hostPort, err := container.MappedPort(ctx, guestPort)
	if err != nil {
		return "", fmt.Errorf("can't resolve host port for the given container ID: %s and guest port: %s; error: %s", container.GetContainerID()[:12], guestPort.Port(), err.Error())
	}
	return hostPort, nil
}

func printPortsSetup(testPackageName string, networkName string, containers AllContainers) {
	hostPortsMessage := fmt.Sprintf("\n>> running test package: '%s' in Docker network name %s; with host port configuration:\n", testPackageName, networkName)
	hostPortsMessage += containers.Redis.SerialisePortSetup()
	hostPortsMessage += containers.Lotus.SerialisePortSetup()
	hostPortsMessage += containers.Register.SerialisePortSetup()
	for _, container := range containers.Gateways {
		hostPortsMessage += container.SerialisePortSetup()
	}
	for _, container := range containers.Providers {
		hostPortsMessage += container.SerialisePortSetup()
	}
	fmt.Println(hostPortsMessage)
}
