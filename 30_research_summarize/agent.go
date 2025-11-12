package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

func main() {
	researchAgent, err := llmagent.New(llmagent.Config{
		Name:        "researchAgent",
		Model:       getModel(),
		Description: "You are a research specialist who details their research with citations",
		Instruction: "You are a specialized research agent. Your only job is to use the google_search tool to find 2 to 3 peieces of relevant information on the given topic and present the findings with citations.",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
		//OutputKey: "research_findings",
	})
	if err != nil {
		log.Fatalf("Failed to create research agent: %v", err)
	}

	summarizerAgent, err := llmagent.New(llmagent.Config{
		Name:        "summarizerAgent",
		Model:       getModel(),
		Instruction: "Read the provided research findings. Create a concise summary as a bulleted list with 3 to 5 key points.",
		//OutputKey:   "final_summary",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:  "researchCoordinator",
		Model: getModel(),
		Instruction: `You are a research coordinator. Your goal is to answer the user's query by orchestrating a workflow.
		1. First you MUST call the 'researchAgent' tool to find relevant information on the topic provided by the user.
		2. Next, after receiving the research findings from researchAgent, you MUST call the 'summarizerAgent' tool with the research findings to create a concise summary .
		3. Finally, present the final summary clearly to the user as your response.`,
		Tools: []tool.Tool{
			agenttool.New(researchAgent, nil),
			agenttool.New(summarizerAgent, nil),
		},
	})
	if err != nil {
		log.Fatalf("Failed to create research agent: %v", err)
	}

	multiRun(agent, researchAgent, summarizerAgent)

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

func multiRun(agent agent.Agent, agents ...agent.Agent) {
	var err error
	//loader, err := services.NewMultiAgentLoader(agent, agents...)
	loader := services.NewSingleAgentLoader(agent)
	if err != nil {
		log.Fatal(err)
	}

	config := &adk.Config{
		AgentLoader: loader,
	}

	l := full.NewLauncher()

	err = l.Execute(context.Background(), config, os.Args[1:])
	if err != nil {
		log.Fatalf("run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
