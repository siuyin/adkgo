package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/loopagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/exitlooptool"
	"google.golang.org/genai"
)

func main() {
	initialWriterAgent, err := llmagent.New(llmagent.Config{
		Name:  "InitialWriterAgent",
		Model: getModel(),
		Instruction: `Based on the user's prompt, write a first draft of a short story
		around 100 to 150 words long.
		Output only the story text, with no introduction or explanation.`,
		OutputKey: "current_story",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	criticAgent, err := llmagent.New(llmagent.Config{
		Name:  "CriticAgent",
		Model: getModel(),
		Instruction: `You are a constructive story critic.
		Review the story below:
		Story: {current_story}

		Evaluate the story's plot, characters and pacing.
		- If the story is well written and complete, you MUST respond with the exact phrase: "APPROVED".
		- Otherwise provide 2 to 3 specific, actionable suggestions for improvement.`,
		OutputKey: "critique",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	exitLoop, err := exitlooptool.New()
	if err != nil {
		log.Fatal(err)
	}

	refinerAgent, err := llmagent.New(llmagent.Config{
		Name:  "RefinerAgent",
		Model: getModel(),
		Instruction: `You are a story refiner. You hava a draft and critique.

		Story Draft: {current_story}
		Critique: {critique}

		Your task is to analyze the critique.
		- If the critique is EXACTLY "APPROVED", you MUST call the exitLoop tool and nothing else.
		- Otherwise rewrite the story draft to incorporate the feedback from the critique.`,
		OutputKey: "current_story",
		Tools:     []tool.Tool{exitLoop},
	})

	storyRefinementLoop, err := loopagent.New(loopagent.Config{
		AgentConfig: agent.Config{
			Name:      "StoryRefinementLoop",
			SubAgents: []agent.Agent{criticAgent, refinerAgent},
		},
		MaxIterations: 3,
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	agent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:      "StoryPipeline",
			SubAgents: []agent.Agent{initialWriterAgent, storyRefinementLoop},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	run(agent)

}

func getModel() model.LLM {
	model, err := gemini.NewModel(context.Background(),
		os.Getenv("MODEL"),
		&genai.ClientConfig{APIKey: os.Getenv("GOOGLE_API_KEY")})
	if err != nil {
		log.Fatalf("could not create model: %v", err)
	}

	return model
}

func run(agent agent.Agent) {
	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(agent),
	}

	l := full.NewLauncher()

	err := l.Execute(context.Background(), config, os.Args[1:])
	if err != nil {
		log.Fatalf("run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
