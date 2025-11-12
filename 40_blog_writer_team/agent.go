package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/genai"
)

func main() {
	outlineAgent, err := llmagent.New(llmagent.Config{
		Name:  "OutlineAgent",
		Model: getModel(),
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
		Model: getModel(),
		Instruction: `Follow this outline strictly: {blog_outline}
		Write a brief 200 to 300 word blog post with an engaging and informative tone.`,
		OutputKey: "blog_draft",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	editorAgent, err := llmagent.New(llmagent.Config{
		Name:  "EditorAgent",
		Model: getModel(),
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
