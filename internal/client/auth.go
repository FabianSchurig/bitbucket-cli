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
// Authentication precedence:
//  1. BITBUCKET_USERNAME + BITBUCKET_API_TOKEN → HTTP Basic Auth (most common)
//  2. BITBUCKET_TOKEN (alone) → Bearer token (OAuth2)
//
// The base URL defaults to https://api.bitbucket.org/2.0 but can be
// overridden with BITBUCKET_BASE_URL (useful for testing).
func NewClient() (*BBClient, error) {
	base := os.Getenv("BITBUCKET_BASE_URL")
	if base == "" {
		base = defaultBaseURL
	}
	c := resty.New().SetBaseURL(base)

	username := os.Getenv("BITBUCKET_USERNAME")
	password := os.Getenv("BITBUCKET_API_TOKEN")
	token := os.Getenv("BITBUCKET_TOKEN")

	switch {
	case username != "" && password != "":
		c.SetBasicAuth(username, password)
	case token != "":
		c.SetAuthToken(token) // Bearer
	default:
		return nil, fmt.Errorf(
			"auth required: set BITBUCKET_USERNAME + BITBUCKET_API_TOKEN, or BITBUCKET_TOKEN",
		)
	}

	return &BBClient{c}, nil
}

// ParseError returns a formatted error from a non-2xx resty response.
func ParseError(resp *resty.Response) error {
	return fmt.Errorf("bitbucket API error %d: %s", resp.StatusCode(), resp.String())
}
