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
//  1. BITBUCKET_USERNAME + BITBUCKET_APP_PASSWORD → HTTP Basic Auth (most common)
//  2. BITBUCKET_TOKEN (alone) → Bearer token (OAuth2)
//
// The base URL defaults to https://api.bitbucket.org/2.0 but can be
// overridden with BITBUCKET_BASE_URL (useful for testing).
func NewClient() (*BBClient, error) {
	return NewClientWithConfig(
		os.Getenv("BITBUCKET_USERNAME"),
		os.Getenv("BITBUCKET_APP_PASSWORD"),
		os.Getenv("BITBUCKET_TOKEN"),
		os.Getenv("BITBUCKET_BASE_URL"),
	)
}

// NewClientWithConfig creates an authenticated Bitbucket API client from
// explicit configuration values. Empty strings are treated as unset.
// This avoids mutating global environment variables.
//
// Authentication precedence:
//  1. username + appPassword → HTTP Basic Auth (legacy app passwords)
//  2. username + token → HTTP Basic Auth (API tokens)
//  3. token alone → Bearer token (OAuth2)
func NewClientWithConfig(username, appPassword, token, baseURL string) (*BBClient, error) {
	base := baseURL
	if base == "" {
		base = defaultBaseURL
	}
	c := resty.New().SetBaseURL(base)

	switch {
	case username != "" && appPassword != "":
		c.SetBasicAuth(username, appPassword)
	case username != "" && token != "":
		c.SetBasicAuth(username, token)
	case token != "":
		c.SetAuthToken(token) // Bearer
	default:
		return nil, fmt.Errorf(
			"auth required: set BITBUCKET_USERNAME + BITBUCKET_TOKEN, or BITBUCKET_TOKEN alone",
		)
	}

	return &BBClient{c}, nil
}

// ParseError returns a formatted error from a non-2xx resty response.
func ParseError(resp *resty.Response) error {
	return fmt.Errorf("bitbucket API error %d: %s", resp.StatusCode(), resp.String())
}
