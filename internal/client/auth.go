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
// Authentication: BITBUCKET_TOKEN is used as a Bearer token.
// BITBUCKET_USERNAME is optional metadata (not used for auth).
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
// Authentication: token is always used as Bearer (OAuth2/API token).
// Username is optional metadata and not used for authentication.
func NewClientWithConfig(username, token, baseURL string) (*BBClient, error) {
	base := baseURL
	if base == "" {
		base = defaultBaseURL
	}
	c := resty.New().SetBaseURL(base)

	switch {
	case token != "":
		c.SetAuthToken(token) // Bearer
	default:
		return nil, fmt.Errorf(
			"auth required: set BITBUCKET_TOKEN (API token or OAuth2 access token)",
		)
	}

	// Store username as metadata (not used for auth).
	_ = username

	return &BBClient{c}, nil
}

// ParseError returns a formatted error from a non-2xx resty response.
func ParseError(resp *resty.Response) error {
	return fmt.Errorf("bitbucket API error %d: %s", resp.StatusCode(), resp.String())
}
