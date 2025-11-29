package main

import (
	"agenttry/mdl"
	"agenttry/runner"
	"log"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/parallelagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
)

func main() {
	techResearcher, err := llmagent.New(llmagent.Config{
		Name:  "TechResearcher",
		Model: mdl.FromEnv(),
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
		Model: mdl.FromEnv(),
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
		Model: mdl.FromEnv(),
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
		Model: mdl.FromEnv(),
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

	runner.Run(agent)

}
