package handlers

import (
	"context"
	"fmt"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/generated"
	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

// MergePRInput holds the validated inputs for the merge pull request command.
type MergePRInput struct {
	Workspace         string
	RepoSlug          string
	PullRequestID     int
	Strategy          string // merge_commit | squash | fast_forward
	Message           string
	CloseSourceBranch bool
}

// MergePullRequest merges the specified pull request.
func MergePullRequest(ctx context.Context, c *client.BBClient, in MergePRInput) error {
	url := fmt.Sprintf(
		"/repositories/%s/%s/pullrequests/%d/merge",
		in.Workspace, in.RepoSlug, in.PullRequestID,
	)

	params := generated.PullrequestMergeParameters{}
	if in.Strategy != "" {
		strategy := generated.PullrequestMergeParametersMergeStrategy(in.Strategy)
		params.MergeStrategy = &strategy
	}
	if in.Message != "" {
		params.Message = &in.Message
	}
	params.CloseSourceBranch = &in.CloseSourceBranch

	var result generated.Pullrequest
	resp, err := c.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post(url)
	if err != nil {
		return fmt.Errorf("merging pull request: %w", err)
	}
	if resp.IsError() {
		return client.ParseError(resp)
	}

	return output.Render(result)
}
