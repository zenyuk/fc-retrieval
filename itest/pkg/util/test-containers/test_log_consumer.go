package test_containers

import (
	"fmt"
	"github.com/ConsenSys/fc-retrieval/itest/pkg/util/test-containers/color"
	"strings"

	tc "github.com/testcontainers/testcontainers-go"
)

type TestLogConsumer struct {
	Messages      []string
	ContainerName string
	Color         string
}

func (g *TestLogConsumer) Accept(l tc.Log) {
	newLogPart := string(l.Content)
	g.Messages = append(g.Messages, newLogPart)
	fmt.Printf("%s%s:%s %s", g.Color, strings.ToUpper(g.ContainerName), color.Reset, newLogPart)
}

/*
type logConsumer struct {
	name  string
	color string
	// The following are used by itest only
	done chan bool
}

func (g *logConsumer) Accept(l tc.Log) {
	log := string(l.Content)
	fmt.Print(g.color, "[", strings.ToUpper(g.name), "]", "\033[0m:", log)
	if g.done != nil {
		if strings.Contains(log, "--- FAIL:") {
			// Tests have falied.
			g.done <- false
		} else if strings.Contains(log, "ok") && strings.Contains(log, "github.com/ConsenSys/fc-retrieval/itest/pkg/") {
			// Tests have all passed.
			g.done <- true
		}
	}
}
*/
