package main

import (
	"agenttry/mdl"
	"agenttry/runner"
	"log"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/geminitool"
)

func main() {
	researchAgent, err := llmagent.New(llmagent.Config{
		Name:        "researchAgent",
		Model:       mdl.FromEnv(),
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
		Model:       mdl.FromEnv(),
		Instruction: "Read the provided research findings. Create a concise summary as a bulleted list with 3 to 5 key points.",
		//OutputKey:   "final_summary",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:  "researchCoordinator",
		Model: mdl.FromEnv(),
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

	runner.Run(agent)

}
