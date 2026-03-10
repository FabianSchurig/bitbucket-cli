package handlers

import (
	"context"
	"fmt"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/generated"
	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

// GetPRInput holds the validated inputs for the get pull request command.
type GetPRInput struct {
	Workspace     string
	RepoSlug      string
	PullRequestID int
}

// GetPullRequest fetches a single pull request by ID.
func GetPullRequest(ctx context.Context, c *client.BBClient, in GetPRInput) error {
	url := fmt.Sprintf(
		"/repositories/%s/%s/pullrequests/%d",
		in.Workspace, in.RepoSlug, in.PullRequestID,
	)

	var pr generated.Pullrequest
	resp, err := c.R().SetContext(ctx).SetResult(&pr).Get(url)
	if err != nil {
		return fmt.Errorf("getting pull request: %w", err)
	}
	if resp.IsError() {
		return client.ParseError(resp)
	}

	return output.Render(pr)
}
