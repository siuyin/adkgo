package main

import (
	"agenttry/mdl"
	"agenttry/runner"
	"log"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/loopagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/exitlooptool"
)

func main() {
	initialWriterAgent, err := llmagent.New(llmagent.Config{
		Name:  "InitialWriterAgent",
		Model: mdl.FromEnv(),
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
		Model: mdl.FromEnv(),
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
		Model: mdl.FromEnv(),
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

	runner.Run(agent)

}
