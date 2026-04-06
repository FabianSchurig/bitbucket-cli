// Package prompts provides MCP prompt handlers (playbooks) for the Bitbucket
// MCP server. Playbooks are embedded markdown files that give LLMs step-by-step
// instructions for common workflows.
package prompts

import (
	"context"
	"embed"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed playbooks/*.md
var playbooks embed.FS

// playbook defines a single MCP prompt backed by an embedded markdown file.
type playbook struct {
	Name        string
	Description string
	File        string // path within the embedded FS
}

var allPlaybooks = []playbook{
	{
		Name:        "bitbucket_pr_reviewer",
		Description: "A step-by-step guide to reviewing a pull request on Bitbucket.",
		File:        "playbooks/pr_reviewer.md",
	},
}

// Register adds all playbook prompts to the MCP server.
func Register(server *mcp.Server) {
	for _, pb := range allPlaybooks {
		pb := pb // capture loop variable
		server.AddPrompt(&mcp.Prompt{
			Name:        pb.Name,
			Description: pb.Description,
		}, promptHandler(pb))
	}
}

func promptHandler(pb playbook) mcp.PromptHandler {
	return func(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		data, err := playbooks.ReadFile(pb.File)
		if err != nil {
			return nil, err
		}
		return &mcp.GetPromptResult{
			Description: pb.Description,
			Messages: []*mcp.PromptMessage{
				{
					Role:    "user",
					Content: &mcp.TextContent{Text: string(data)},
				},
			},
		}, nil
	}
}
