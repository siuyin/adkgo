package main

import (
	"agenttry/mdl"
	"agenttry/runner"
	"log"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
)

func main() {
	outlineAgent, err := llmagent.New(llmagent.Config{
		Name:  "OutlineAgent",
		Model: mdl.FromEnv(),
		Instruction: `Create a blog outline for the given topic with:
		1. A catchy headline
		2. An introduction hook
		3. 3 to 5 main sections with 2 to 3 bullet points each.
		4. A concluding thought. `,
		OutputKey: "blog_outline",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	writerAgent, err := llmagent.New(llmagent.Config{
		Name:  "WriterAgent",
		Model: mdl.FromEnv(),
		Instruction: `Follow this outline strictly: {blog_outline}
		Write a brief 200 to 300 word blog post with an engaging and informative tone.`,
		OutputKey: "blog_draft",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	editorAgent, err := llmagent.New(llmagent.Config{
		Name:  "EditorAgent",
		Model: mdl.FromEnv(),
		Instruction: `Edit this draft: {blog_draft}
		Your task is to polish the text by fixing any grammatical errors, improve the flow and sentence structure, and enhance overall clarity. `,
		OutputKey: "final_blog",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	agent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{Name: "sequentialAgent", SubAgents: []agent.Agent{outlineAgent, writerAgent, editorAgent}},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	runner.Run(agent)

}
