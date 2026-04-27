// Package client provides an authenticated HTTP client for the Bitbucket API.
package client

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

const defaultBaseURL = "https://api.bitbucket.org/2.0"

// BBClient wraps a resty.Client configured for the Bitbucket API.
//
// It carries two independent sets of credentials:
//
//   - Username + Token: HTTP Basic Auth, used for the public REST API
//     (api.bitbucket.org/2.0).
//   - CSRFToken + CloudSessionToken: cookie-based auth used by Bitbucket's
//     undocumented internal API (https://bitbucket.org/!api/internal/...).
//     The internal API does NOT accept HTTP Basic Auth — it requires the same
//     csrftoken + cloud.session.token cookies and X-CSRFToken header that the
//     Bitbucket web UI sends. The dispatcher inspects the request URL and
//     applies the appropriate credentials per request.
type BBClient struct {
	*resty.Client
	Username          string
	Token             string
	CSRFToken         string
	CloudSessionToken string
}

// NewClient creates an authenticated Bitbucket API client from environment
// variables.
//
// Public REST API auth (one of):
//   - BITBUCKET_USERNAME + BITBUCKET_TOKEN → HTTP Basic Auth
//   - BITBUCKET_TOKEN alone               → HTTP Basic Auth with "x-token-auth"
//
// Internal API auth (both required for /!api/internal/ endpoints):
//   - BITBUCKET_CSRF_TOKEN
//   - BITBUCKET_CLOUD_SESSION_TOKEN
//
// At least one of the two auth modes must be configured.
//
// The base URL defaults to https://api.bitbucket.org/2.0 but can be
// overridden with BITBUCKET_BASE_URL (useful for testing).
func NewClient() (*BBClient, error) {
	return NewClientWithConfig(
		os.Getenv("BITBUCKET_USERNAME"),
		os.Getenv("BITBUCKET_TOKEN"),
		os.Getenv("BITBUCKET_BASE_URL"),
		os.Getenv("BITBUCKET_CSRF_TOKEN"),
		os.Getenv("BITBUCKET_CLOUD_SESSION_TOKEN"),
	)
}

// NewClientWithConfig creates an authenticated Bitbucket API client from
// explicit configuration values. Empty strings are treated as unset.
// This avoids mutating global environment variables.
//
// Authentication precedence (per request, decided by the dispatcher):
//   - URL contains "/!api/internal/": csrfToken + cloudSessionToken cookies
//     and X-CSRFToken header. Basic Auth is suppressed.
//   - Otherwise: HTTP Basic Auth using username + token (or "x-token-auth" +
//     token when username is empty, for workspace/repository access tokens).
func NewClientWithConfig(username, token, baseURL, csrfToken, cloudSessionToken string) (*BBClient, error) {
	base := baseURL
	if base == "" {
		base = defaultBaseURL
	}
	c := resty.New().SetBaseURL(base)

	hasBasic := token != ""
	hasInternal := csrfToken != "" && cloudSessionToken != ""
	if !hasBasic && !hasInternal {
		return nil, fmt.Errorf(
			"auth required: set BITBUCKET_TOKEN for the public API, " +
				"or set BITBUCKET_CSRF_TOKEN and BITBUCKET_CLOUD_SESSION_TOKEN " +
				"to access the internal API (basic auth is not supported there)",
		)
	}

	// Pre-apply Basic Auth on the resty client so callers that bypass the
	// dispatcher (e.g. ad-hoc tooling) still get authenticated for the public
	// API. The dispatcher overrides this per-request for internal URLs.
	if hasBasic {
		authUser := username
		if authUser == "" {
			authUser = "x-token-auth"
		}
		c.SetBasicAuth(authUser, token)
	}

	return &BBClient{
		Client:            c,
		Username:          username,
		Token:             token,
		CSRFToken:         csrfToken,
		CloudSessionToken: cloudSessionToken,
	}, nil
}

// ParseError returns a formatted error from a non-2xx resty response.
func ParseError(resp *resty.Response) error {
	return fmt.Errorf("bitbucket API error %d: %s", resp.StatusCode(), resp.String())
}
