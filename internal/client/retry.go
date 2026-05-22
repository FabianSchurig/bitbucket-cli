package client

import (
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

// ConfigureRetry sets up resty's built-in retry mechanism on the given client.
// It uses exponential backoff and retries only on transient HTTP errors.
// This is safe to call on any resty.Client and benefits all consumers
// (CLI, MCP server, Terraform provider) uniformly.
func ConfigureRetry(c *resty.Client) {
	c.SetRetryCount(retryMaxAttempts).
		SetRetryWaitTime(retryInitialWait).
		SetRetryMaxWaitTime(retryMaxWait).
		AddRetryCondition(func(resp *resty.Response, err error) bool {
			if err != nil {
				// Network-level errors (timeouts, connection refused) are retryable.
				return true
			}
			return retryableStatusCodes[resp.StatusCode()]
		})
}
