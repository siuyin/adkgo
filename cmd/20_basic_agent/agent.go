package main

import (
	"agenttry/mdl"
	"agenttry/runner"
	"log"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
)

func main() {
	agent, err := llmagent.New(llmagent.Config{
		Name:        "helpful_assistant",
		Model:       mdl.FromEnv(),
		Description: "A simple agent that can answer general questions.",
		Instruction: "You are a helpful assistant. Use google search for current info or if unsure.",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	runner.Run(agent)

}
