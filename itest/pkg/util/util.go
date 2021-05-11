package util

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go/wait"
	tc "github.com/wcgcyx/testcontainers-go"
)

const ColorRed = "\033[31m"
const ColorGreen = "\033[32m"
const ColorYellow = "\033[33m"
const ColorBlue = "\033[34m"
const ColorPurple = "\033[35m"
const ColorCyan = "\033[36m"
const ColorWhite = "\033[37m" // Used by redis
const ColorGray = "\033[90m"  // Used by lotus
const ColorBrightRed = "\033[91m"
const ColorBrightGreen = "\033[91m"
const ColorBrightYellow = "\033[91m"
const ColorBrightBlue = "\033[91m"
const ColorBrightPurple = "\033[91m"
const ColorBrightCyan = "\033[91m"
const ColorBrightWhite = "\033[91m"

func StartContainers() (string, error) {
	composeFilePaths := []string{"../../docker-compose.yml"}
	identifier := strings.ToLower(uuid.New().String())

	compose := tc.NewLocalDockerCompose(composeFilePaths, identifier)
	execError := compose.
		WithCommand([]string{"up", "-d"}).
		Invoke()

	err := execError.Error
	if err != nil {
		return "", fmt.Errorf("could not run compose file: %v - %v", composeFilePaths, err)
	}
	return compose.Identifier, nil
}

func StopContainers(composeID string) error {
	composeFilePaths := []string{"../../docker-compose.yml"}

	compose := tc.NewLocalDockerCompose(composeFilePaths, composeID)
	execError := compose.Down()
	err := execError.Error
	if err != nil {
		return fmt.Errorf("could not stop compose file: %v - %v", composeFilePaths, err)
	}
	return nil
}

// CleanContainers clean the containers
func CleanContainers() {
	// Stop all running containers
	cmd := exec.Command("docker", "ps", "-q")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	containers := strings.Split(string(stdout), "\n")
	command := append([]string{"stop"}, containers...)
	cmd = exec.Command("docker", command...)
	stdout, _ = cmd.Output()
	fmt.Printf("Stopped containers: \n%v\n", string(stdout))

	// Remove all stopped containers
	cmd = exec.Command("docker", "ps", "-a", "-q")
	stdout, err = cmd.Output()
	if err != nil {
		panic(err)
	}
	containers = strings.Split(string(stdout), "\n")
	command = append([]string{"rm"}, containers...)
	cmd = exec.Command("docker", command...)
	stdout, _ = cmd.Output()
	fmt.Printf("Removed containers: \n%v\n", string(stdout))
}

// GetCurrentBranch gets the current branch of this repo
func GetCurrentBranch() string {
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	tag := string(stdout[:len(stdout)-1])
	return tag
}

// CreateNetwork creates a network
func CreateNetwork(ctx context.Context, network string) *tc.Network {
	// First remove the network if existed
	cmd := exec.Command("docker", "network", "rm", network)
	cmd.Output()

	net, err := tc.GenericNetwork(ctx, tc.GenericNetworkRequest{
		NetworkRequest: tc.NetworkRequest{
			Name:           network,
			CheckDuplicate: true,
		},
	})
	if err != nil {
		panic(err)
	}
	return &net
}

// GetEnvMap gets the env map from a given env file
func GetEnvMap(envFile string) map[string]string {
	env := make(map[string]string)
	file, err := os.Open(envFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

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
	return env
}

// Start lotus
func StartLotus(ctx context.Context, network string, verbose bool) *tc.Container {
	// Start lotus
	req := tc.ContainerRequest{
		Image:          "consensys/lotus-full-node:latest",
		Name:           "lotus",
		Networks:       []string{network},
		NetworkMode:    container.NetworkMode(network),
		NetworkAliases: map[string][]string{network: {"lotus"}},
		WaitingFor:     wait.ForLog("mined new block"),
	}
	lotusC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	if verbose {
		g := &logConsumer{name: fmt.Sprintf("lotus"), color: ColorGray}
		err = lotusC.StartLogProducer(ctx)
		if err != nil {
			panic(err)
		}
		lotusC.FollowOutput(g)
	}
	return &lotusC
}

// Start redis used by register
func StartRedis(ctx context.Context, network string) *tc.Container {
	// Start redis
	req := tc.ContainerRequest{
		Image:          "redis:alpine",
		Name:           "redis",
		Cmd:            []string{"redis-server", "--requirepass", "xxxx"},
		Networks:       []string{network},
		NetworkMode:    container.NetworkMode(network),
		NetworkAliases: map[string][]string{network: {"redis"}},
		WaitingFor:     wait.ForLog("Ready to accept connections"),
	}
	redisC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	g := &logConsumer{name: fmt.Sprintf("redis"), color: ColorWhite}
	err = redisC.StartLogProducer(ctx)
	if err != nil {
		panic(err)
	}
	redisC.FollowOutput(g)
	return &redisC
}

// Start the register
func StartRegister(ctx context.Context, tag string, network string, color string, env map[string]string) *tc.Container {
	// Start a register container
	req := tc.ContainerRequest{
		Image:          fmt.Sprintf("consensys/fc-retrieval-register:develop-%s", tag),
		Name:           "register",
		Networks:       []string{network},
		Env:            env,
		NetworkMode:    container.NetworkMode(network),
		NetworkAliases: map[string][]string{network: {"register"}},
		WaitingFor:     wait.ForLog("Serving register at"),
	}
	registerC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	g := &logConsumer{name: fmt.Sprintf("register"), color: color}
	err = registerC.StartLogProducer(ctx)
	if err != nil {
		panic(err)
	}
	registerC.FollowOutput(g)
	return &registerC
}

// Start a gateway of specific id, tag, network, log color and env
func StartGateway(ctx context.Context, id string, tag string, network string, color string, env map[string]string) *tc.Container {
	// Start a gateway container
	req := tc.ContainerRequest{
		Image:          fmt.Sprintf("consensys/fc-retrieval-gateway:develop-%s", tag),
		Name:           id,
		Networks:       []string{network},
		Env:            env,
		NetworkMode:    container.NetworkMode(network),
		NetworkAliases: map[string][]string{network: {id}},
		WaitingFor:     wait.ForLog("Filecoin Gateway Start-up Complete"),
	}
	gatewayC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	g := &logConsumer{name: id, color: color}
	err = gatewayC.StartLogProducer(ctx)
	if err != nil {
		panic(err)
	}
	gatewayC.FollowOutput(g)

	return &gatewayC
}

// Start a provider of specific id, tag, network, log color and env
func StartProvider(ctx context.Context, id string, tag string, network string, color string, env map[string]string) *tc.Container {
	// Start a provider container
	req := tc.ContainerRequest{
		Image:          fmt.Sprintf("consensys/fc-retrieval-provider:develop-%s", tag),
		Name:           id,
		Networks:       []string{network},
		Env:            env,
		NetworkMode:    container.NetworkMode(network),
		NetworkAliases: map[string][]string{network: {id}},
		WaitingFor:     wait.ForLog("Filecoin Provider Start-up Complete"),
	}
	providerC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	g := &logConsumer{name: id, color: color}
	err = providerC.StartLogProducer(ctx)
	if err != nil {
		panic(err)
	}
	providerC.FollowOutput(g)
	return &providerC
}

// Start the itest
func StartItest(ctx context.Context, tag string, network string, color string, testDir string, done chan bool) *tc.Container {
	// Start a itest container
	req := tc.ContainerRequest{
		Image:          fmt.Sprintf("consensys/fc-retrieval-itest:develop-%s", tag),
		Name:           "itest",
		Networks:       []string{network},
		Env:            map[string]string{"ITEST_CALLING_FROM_CONTAINER": "yes"},
		NetworkMode:    container.NetworkMode(network),
		NetworkAliases: map[string][]string{network: {"itest"}},
		Cmd:            []string{"go", "test", "-v", "--count=1", testDir},
	}
	itestC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	g := &logConsumer{name: fmt.Sprintf("itest"), color: color, done: done}
	err = itestC.StartLogProducer(ctx)
	if err != nil {
		panic(err)
	}
	itestC.FollowOutput(g)
	return &itestC
}

type logConsumer struct {
	name  string
	color string
	// The following are used by itest only
	done chan bool
}

func (g *logConsumer) Accept(l tc.Log) {
	log := string(l.Content)
	fmt.Print("{", string(g.color), g.name, "}: ", string("\033[0m"), log)
	if g.done != nil {
		if strings.Contains(log, "--- FAIL:") {
			// Tests have falied.
			g.done <- false
		} else if strings.Contains(log, "ok") && strings.Contains(log, "github.com/ConsenSys/fc-retrieval-itest/pkg/") {
			// Tests have all passed.
			g.done <- true
		}
	}
}
