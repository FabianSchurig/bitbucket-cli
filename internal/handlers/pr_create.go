package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/generated"
	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

// CreatePRInput holds the validated inputs for the create pull request command.
type CreatePRInput struct {
	Workspace          string
	RepoSlug           string
	Title              string
	Description        string
	SourceBranch       string
	DestinationBranch  string
	ReviewerUUIDs      []string // comma-separated UUIDs expanded into reviewers array
	CloseSourceBranch  bool
	Body               string // raw JSON body escape hatch; overrides all other fields if set
}

// CreatePullRequest creates a new pull request.
// If in.Body is non-empty, it is used as the raw JSON body (power-user escape hatch).
// Otherwise, the individual flag values are assembled into the request body.
func CreatePullRequest(ctx context.Context, c *client.BBClient, in CreatePRInput) error {
	url := fmt.Sprintf("/repositories/%s/%s/pullrequests", in.Workspace, in.RepoSlug)

	var pr generated.Pullrequest

	if in.Body != "" {
		// Raw JSON body escape hatch
		if err := json.Unmarshal([]byte(in.Body), &pr); err != nil {
			return fmt.Errorf("parsing --body JSON: %w", err)
		}
	} else {
		// Assemble from individual flags
		pr.Title = &in.Title
		if in.Description != "" {
			pr.Description = &in.Description
		}
		if in.SourceBranch != "" {
			branchName := in.SourceBranch
			pr.Source = &generated.PullrequestEndpoint{
				Branch: &struct {
					Name *string `json:"name,omitempty"`
				}{Name: &branchName},
			}
		}
		if in.DestinationBranch != "" {
			branchName := in.DestinationBranch
			pr.Destination = &generated.PullrequestEndpoint{
				Branch: &struct {
					Name *string `json:"name,omitempty"`
				}{Name: &branchName},
			}
		}
		if len(in.ReviewerUUIDs) > 0 {
			var reviewers []generated.Account
			for _, raw := range in.ReviewerUUIDs {
				for _, uuid := range strings.Split(raw, ",") {
					uuid = strings.TrimSpace(uuid)
					if uuid != "" {
						u := uuid
						reviewers = append(reviewers, generated.Account{Uuid: &u})
					}
				}
			}
			pr.Reviewers = &reviewers
		}
	}

	var result generated.Pullrequest
	resp, err := c.R().
		SetContext(ctx).
		SetBody(pr).
		SetResult(&result).
		Post(url)
	if err != nil {
		return fmt.Errorf("creating pull request: %w", err)
	}
	if resp.IsError() {
		return client.ParseError(resp)
	}

	return output.Render(result)
}
