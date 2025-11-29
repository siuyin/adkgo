package runner

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
)

// Run runs the agent.
func Run(agt agent.Agent) {
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(agt),
	}

	l := full.NewLauncher()

	err := l.Execute(context.Background(), config, os.Args[1:])
	if err != nil {
		log.Fatalf("run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
