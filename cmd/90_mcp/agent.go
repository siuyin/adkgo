package main

import (
	"agenttry/mdl"
	"agenttry/runner"
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/oauth2"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
)

func main() {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_PAT")})

	githubMCPServer, err := mcptoolset.New(mcptoolset.Config{
		Transport: &mcp.StreamableClientTransport{
			Endpoint:   "https://api.githubcopilot.com/mcp/",
			HTTPClient: oauth2.NewClient(context.Background(), ts),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = githubMCPServer

	mcpImageServer, err := mcptoolset.New(mcptoolset.Config{
		Transport:  &mcp.CommandTransport{Command: exec.Command("npx", "-y", "@modelcontextprotocol/server-everything")},
		ToolFilter: tool.StringPredicate([]string{"getTinyImage", "echo"}),
	})
	if err != nil {
		log.Fatal(err)
	}

	mcpAgent, err := llmagent.New(llmagent.Config{
		Name:  "ImageAgent",
		Model: mdl.FromEnv(),
		Instruction: `Use the MCP tool from the toolset to generate images or to parrot/echo  user queries
		when calling the getTinyImage tool issue the 'includeImage' argument.
		The getTinyImage tool produces base-64 encoded data. When asked to generate images, output the raw data verbatim and state that is base-64 encoded.`,
		Toolsets: []tool.Toolset{mcpImageServer, githubMCPServer}, // To list tools: use a prompt like: "What tools do you have?"
	})
	if err != nil {
		log.Fatal(err)
	}

	runner.Run(mcpAgent)

	//sess := session.InMemoryService()
	//sess.Create(context.Background(), &session.CreateRequest{AppName: "myApp", UserID: "myUser", SessionID: "sess-01"})
	//r, err := runner.New(runner.Config{
	//	AppName:        "myApp",
	//	Agent:          mcpAgent,
	//	SessionService: sess,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}

	//resp := r.Run(context.Background(), "myUser", "sess-01", genai.NewContentFromText("provide a sample tiny image", "user"), agent.RunConfig{SaveInputBlobsAsArtifacts: true})
	//for s, e := range resp {
	//	if e != nil {
	//		log.Println(err)
	//		break
	//	}
	//	//fmt.Printf("%#v", s.Content)
	//	for _, p := range s.Content.Parts {
	//		fmt.Printf("%#v\n\n", p)
	//	}
	//}

}
