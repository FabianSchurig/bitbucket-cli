package client

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// retryMaxAttempts is the number of retries after the initial request.
	retryMaxAttempts = 3

	// retryInitialWait is the base wait time before the first retry.
	retryInitialWait = 500 * time.Millisecond

	// retryMaxWait is the upper bound for exponential backoff wait time.
	retryMaxWait = 5 * time.Second
)

// retryableStatusCodes defines HTTP status codes that indicate transient
// failures worth retrying. This includes:
//   - 429: Too Many Requests (rate limiting)
//   - 502: Bad Gateway
//   - 503: Service Unavailable
//   - 504: Gateway Timeout
var retryableStatusCodes = map[int]bool{
	429: true,
	502: true,
	503: true,
	504: true,
}

// idempotentMethods defines HTTP methods that are safe to retry without
// risking duplicate side effects.
var idempotentMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodPut:     true,
	http.MethodDelete:  true,
	http.MethodOptions: true,
}

// ConfigureRetry sets up resty's built-in retry mechanism on the given client.
// It uses exponential backoff and retries only on transient HTTP errors for
// idempotent methods (GET, HEAD, PUT, DELETE, OPTIONS). Non-idempotent methods
// like POST and PATCH are never retried to avoid duplicate side effects.
// Context cancellation and deadline errors are also not retried.
// This is safe to call on any resty.Client and benefits all consumers
// (CLI, MCP server, Terraform provider) uniformly.
func ConfigureRetry(c *resty.Client) {
	c.SetRetryCount(retryMaxAttempts).
		SetRetryWaitTime(retryInitialWait).
		SetRetryMaxWaitTime(retryMaxWait).
		AddRetryCondition(retryCondition)
}

// retryCondition determines whether a failed request should be retried.
func retryCondition(resp *resty.Response, err error) bool {
	// Never retry context cancellation or deadline exceeded errors.
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return false
		}
	}

	// Determine the HTTP method from the response or request.
	method := ""
	if resp != nil && resp.Request != nil && resp.Request.RawRequest != nil {
		method = resp.Request.RawRequest.Method
	}

	// Only retry idempotent methods to avoid duplicate side effects.
	if !idempotentMethods[method] {
		return false
	}

	if err != nil {
		// Network-level errors (timeouts, connection refused) are retryable.
		return true
	}
	return retryableStatusCodes[resp.StatusCode()]
}
