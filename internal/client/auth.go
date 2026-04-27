// Package client provides an authenticated HTTP client for the Bitbucket API.
package client

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	defaultBaseURL = "https://api.bitbucket.org/2.0"

	// internalAPIMarker is the URL substring that identifies Bitbucket's
	// undocumented internal API. Kept here (and mirrored in the dispatcher)
	// so the client-level pre-request hook can recognise internal URLs
	// without depending on the handlers package.
	internalAPIMarker = "/!api/internal/"
)

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

	// Defence in depth: even though the dispatcher applies Basic Auth
	// per-request (and not at the client level), this hook guarantees that
	// no Authorization header ever reaches an internal-API endpoint, even
	// if a future caller wires Basic Auth onto the underlying resty client
	// directly. The internal endpoint returns 401 when both cookies and an
	// Authorization header are present.
	c.SetPreRequestHook(func(_ *resty.Client, r *http.Request) error {
		if r != nil && r.URL != nil && strings.Contains(r.URL.Path, internalAPIMarker) {
			r.Header.Del("Authorization")
		}
		return nil
	})

	hasBasic := token != ""
	hasInternal := csrfToken != "" && cloudSessionToken != ""
	if !hasBasic && !hasInternal {
		return nil, fmt.Errorf(
			"auth required: set BITBUCKET_TOKEN for the public API, " +
				"or set BITBUCKET_CSRF_TOKEN and BITBUCKET_CLOUD_SESSION_TOKEN " +
				"to access the internal API (basic auth is not supported there)",
		)
	}

	// Authentication is selected per request by the dispatcher
	// (handlers.executeRequest) based on the URL: public REST API → Basic
	// Auth using Username/Token; internal API (/!api/internal/) → cookies
	// + X-CSRFToken. We deliberately do NOT call resty's client-level
	// SetBasicAuth: doing so would force the dispatcher into fragile
	// header-deletion tricks to suppress Basic Auth on internal-API
	// requests (resty re-injects the Authorization header from its
	// client-level UserInfo during request execution).

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
