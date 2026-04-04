// Package client provides an authenticated HTTP client for the Bitbucket API.
package client

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

const defaultBaseURL = "https://api.bitbucket.org/2.0"

// BBClient wraps a resty.Client configured for the Bitbucket API.
type BBClient struct {
	*resty.Client
}

// NewClient creates an authenticated Bitbucket API client.
//
// Authentication: BITBUCKET_USERNAME + BITBUCKET_TOKEN → HTTP Basic Auth.
// If only BITBUCKET_TOKEN is set, "x-token-auth" is used as the username
// (standard for Bitbucket workspace/repository access tokens).
//
// The base URL defaults to https://api.bitbucket.org/2.0 but can be
// overridden with BITBUCKET_BASE_URL (useful for testing).
func NewClient() (*BBClient, error) {
	return NewClientWithConfig(
		os.Getenv("BITBUCKET_USERNAME"),
		os.Getenv("BITBUCKET_TOKEN"),
		os.Getenv("BITBUCKET_BASE_URL"),
	)
}

// NewClientWithConfig creates an authenticated Bitbucket API client from
// explicit configuration values. Empty strings are treated as unset.
// This avoids mutating global environment variables.
//
// Authentication precedence:
//   - username + token → HTTP Basic Auth (works for app passwords and personal tokens)
//   - token alone → HTTP Basic Auth with "x-token-auth" as the username
//     (standard method for Bitbucket workspace and repository access tokens)
func NewClientWithConfig(username, token, baseURL string) (*BBClient, error) {
	base := baseURL
	if base == "" {
		base = defaultBaseURL
	}
	c := resty.New().SetBaseURL(base)

	if token == "" {
		return nil, fmt.Errorf(
			"auth required: set BITBUCKET_TOKEN",
		)
	}

	// When a username is provided, use it directly (app passwords, personal tokens).
	// Otherwise fall back to "x-token-auth" (workspace/repository access tokens).
	authUser := username
	if authUser == "" {
		authUser = "x-token-auth"
	}
	c.SetBasicAuth(authUser, token)

	return &BBClient{c}, nil
}

// ParseError returns a formatted error from a non-2xx resty response.
func ParseError(resp *resty.Response) error {
	return fmt.Errorf("bitbucket API error %d: %s", resp.StatusCode(), resp.String())
}
