package util

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
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

// CleanContainers clean the containers
func CleanContainers(network string) {
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

	// Remove network
	cmd = exec.Command("docker", "network", "rm", network)
	cmd.Output()
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

// GetImageTag gets the image tag of a given repo and tag
func GetImageTag(repo, tag string) string {
	localImageTag := fmt.Sprintf("%v:develop-%v", repo, tag)
	localImageMain := fmt.Sprintf("%v:develop-main", repo)
	remoteImage := fmt.Sprintf("%v:dev", repo)

	cmd := exec.Command("docker", "image", "inspect", localImageTag)
	_, err := cmd.Output()
	if err == nil {
		return localImageTag
	}

	cmd = exec.Command("docker", "image", "inspect", localImageMain)
	_, err = cmd.Output()
	if err == nil {
		return localImageMain
	}

	return remoteImage
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

// StartLotus starts Lotus local development network, two services: miner and daemon in one container
func StartLotus(ctx context.Context, network string, verbose bool) {
	// Start lotus
	req := tc.ContainerRequest{
		Image:          "consensys/lotus-full-node:latest",
		Name:           "lotus",
		Cmd:            []string{"./start-lotus-full-node.sh"},
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
}

// Start redis used by register
func StartRedis(ctx context.Context, network string, verbose bool) {
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
	if verbose {
		g := &logConsumer{name: fmt.Sprintf("redis"), color: ColorWhite}
		err = redisC.StartLogProducer(ctx)
		if err != nil {
			panic(err)
		}
		redisC.FollowOutput(g)
	}
}

// Start the register
func StartRegister(ctx context.Context, tag string, network string, color string, env map[string]string, verbose bool) {
	// Start a register container
	req := tc.ContainerRequest{
		Image:          GetImageTag("consensys/fc-retrieval-register", tag),
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
	if verbose {
		g := &logConsumer{name: fmt.Sprintf("register"), color: color}
		err = registerC.StartLogProducer(ctx)
		if err != nil {
			panic(err)
		}
		registerC.FollowOutput(g)
	}
}

// Start a gateway of specific id, tag, network, log color and env
func StartGateway(ctx context.Context, id string, tag string, network string, color string, env map[string]string, verbose bool) {
	// Start a gateway container
	req := tc.ContainerRequest{
		Image:          GetImageTag("consensys/fc-retrieval-gateway", tag),
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
	if verbose {
		g := &logConsumer{name: id, color: color}
		err = gatewayC.StartLogProducer(ctx)
		if err != nil {
			panic(err)
		}
		gatewayC.FollowOutput(g)
	}
}

// Start a provider of specific id, tag, network, log color and env
func StartProvider(ctx context.Context, id string, tag string, network string, color string, env map[string]string, verbose bool) {
	// Start a provider container
	req := tc.ContainerRequest{
		Image:          GetImageTag("consensys/fc-retrieval-provider", tag),
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
	if verbose {
		g := &logConsumer{name: id, color: color}
		err = providerC.StartLogProducer(ctx)
		if err != nil {
			panic(err)
		}
		providerC.FollowOutput(g)
	}
}

// Start the itest, must only be called in host
func StartItest(ctx context.Context, tag string, network string, color string, done chan bool, verbose bool) {
	// Start a itest container
	// Mount testdir
	absPath, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	req := tc.ContainerRequest{
		Image:          GetImageTag("consensys/fc-retrieval-itest", tag),
		Name:           "itest",
		Networks:       []string{network},
		Env:            map[string]string{"ITEST_CALLING_FROM_CONTAINER": "yes"},
		NetworkMode:    container.NetworkMode(network),
		NetworkAliases: map[string][]string{network: {"itest"}},
		BindMounts:     map[string]string{absPath: "/go/src/github.com/ConsenSys/fc-retrieval-itest/pkg/temp/"},
		Cmd:            []string{"go", "test", "-v", "--count=1", "/go/src/github.com/ConsenSys/fc-retrieval-itest/pkg/temp/"},
	}
	itestC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	if verbose {
		g := &logConsumer{name: fmt.Sprintf("itest"), color: color, done: done}
		err = itestC.StartLogProducer(ctx)
		if err != nil {
			panic(err)
		}
		itestC.FollowOutput(g)
	}
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
