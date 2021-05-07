package test_containers

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
)

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
