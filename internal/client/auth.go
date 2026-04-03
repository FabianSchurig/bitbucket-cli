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
// Authentication: BITBUCKET_TOKEN is used with HTTP Basic Auth
// (x-token-auth:{token}), the standard method for Bitbucket
// workspace and repository access tokens.
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
// Authentication: uses the API token with HTTP Basic Auth (x-token-auth).
// This is the standard method for Bitbucket workspace and repository access tokens.
func NewClientWithConfig(username, token, baseURL string) (*BBClient, error) {
	base := baseURL
	if base == "" {
		base = defaultBaseURL
	}
	c := resty.New().SetBaseURL(base)

	if token == "" {
		return nil, fmt.Errorf(
			"auth required: set BITBUCKET_TOKEN (workspace or repository access token)",
		)
	}
	// Bitbucket access tokens authenticate via Basic Auth with the
	// fixed username "x-token-auth" and the token as the password.
	c.SetBasicAuth("x-token-auth", token)

	// Username is stored for identification purposes only.
	_ = username

	return &BBClient{c}, nil
}

// ParseError returns a formatted error from a non-2xx resty response.
func ParseError(resp *resty.Response) error {
	return fmt.Errorf("bitbucket API error %d: %s", resp.StatusCode(), resp.String())
}
