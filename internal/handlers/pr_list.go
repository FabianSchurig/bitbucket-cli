// Package handlers implements the Bitbucket API dispatch layer for each CLI command.
package handlers

import (
	"context"
	"fmt"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/generated"
	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

// ListPRsInput holds the validated inputs for the list pull requests command.
type ListPRsInput struct {
	Workspace string
	RepoSlug  string
	State     string // empty = all states
	All       bool   // traverse all pages
}

// ListPullRequests fetches pull requests for a repository.
// When in.All is true, it follows the "next" cursor until all pages are fetched.
func ListPullRequests(ctx context.Context, c *client.BBClient, in ListPRsInput) error {
	var allPRs []generated.Pullrequest

	nextURL := fmt.Sprintf("/repositories/%s/%s/pullrequests", in.Workspace, in.RepoSlug)

	for nextURL != "" {
		req := c.R().SetContext(ctx)
		if in.State != "" {
			req = req.SetQueryParam("state", in.State)
		}

		var page generated.PaginatedPullrequests
		resp, err := req.SetResult(&page).Get(nextURL)
		if err != nil {
			return fmt.Errorf("listing pull requests: %w", err)
		}
		if resp.IsError() {
			return client.ParseError(resp)
		}

		if page.Values != nil {
			allPRs = append(allPRs, *page.Values...)
		}

		// Follow Bitbucket's cursor-based "next" pagination link.
		if in.All && page.Next != nil && *page.Next != "" {
			nextURL = *page.Next // absolute URL — resty handles this correctly
		} else {
			break
		}
	}

	return output.Render(allPRs)
}
