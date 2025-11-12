package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/parallelagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

func main() {
	techResearcher, err := llmagent.New(llmagent.Config{
		Name:  "TechResearcher",
		Model: getModel(),
		Instruction: `Research the latest AI/ML trends.
		Include 3 key developments,
		the main companies involved,
		and the development's potential impact.
		Keep the report concise and within 100 words.`,
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
		OutputKey: "tech_research",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	healthReseacher, err := llmagent.New(llmagent.Config{
		Name:  "HealthResearcher",
		Model: getModel(),
		Instruction: `Research recent medical breakthroughs.
		Include 3 significant advances, their practical applications,
		and estimated timelines.
		Keep the report concise and within 100 words.`,
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
		OutputKey: "health_research",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	financeResearcher, err := llmagent.New(llmagent.Config{
		Name:  "FinanceResearcher",
		Model: getModel(),
		Instruction: `Research current fintech trends.
		Include 3 key trends, their market implications and their future outlook.
		Keep the report concise and within 100 words.`,
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
		OutputKey: "finance_research",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	aggregatorAgent, err := llmagent.New(llmagent.Config{
		Name:  "AggregatorAgent",
		Model: getModel(),
		Instruction: `Combine these three research findings into an executive summary.

		**Technology Trends:**
		{tech_research}

		**Health Breakthroughs:**
		{health_research}

		**Finance Innovations:**
		{finance_research}

		Your summary should highlight common themes, surprising or interesting connections and the key takeaways from all 3 reports.
		The final summary should be about 200 words.`,
		OutputKey: "executive_summary",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	parallelResearchTeam, err := parallelagent.New(parallelagent.Config{
		AgentConfig: agent.Config{Name: "ParallelResearchTeam",
			SubAgents: []agent.Agent{techResearcher, healthReseacher, financeResearcher}},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	agent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:      "ResearchSystem",
			SubAgents: []agent.Agent{parallelResearchTeam, aggregatorAgent},
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
