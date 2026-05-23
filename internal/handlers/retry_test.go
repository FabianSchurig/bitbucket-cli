package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/handlers"
	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

// newRetryTestClient creates a BBClient with retry configured, pointing at the test server.
// Retry wait times are overridden to very small durations so tests run fast.
func newRetryTestClient(t *testing.T, serverURL string) *client.BBClient {
	t.Helper()
	r := resty.New().SetBaseURL(serverURL)
	client.ConfigureRetry(r)
	r.SetRetryWaitTime(1 * time.Millisecond).SetRetryMaxWaitTime(5 * time.Millisecond)
	return &client.BBClient{Client: r, Token: "test-token"}
}

func TestRetry_429_SucceedsOnSecondAttempt(t *testing.T) {
	output.Format = "json"
	var attempts int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"message":"rate limited"}}`))
			return
		}
		w.Header().Set(headerContentType, contentTypeJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1}`))
	}))
	defer srv.Close()

	c := newRetryTestClient(t, srv.URL)
	_, err := handlers.DispatchRaw(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: "/test",
	})
	if err != nil {
		t.Fatalf("expected success after retry, got: %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Errorf("expected 2 attempts, got %d", got)
	}
}

func TestRetry_503_SucceedsOnThirdAttempt(t *testing.T) {
	output.Format = "json"
	var attempts int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error":{"message":"service unavailable"}}`))
			return
		}
		w.Header().Set(headerContentType, contentTypeJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1}`))
	}))
	defer srv.Close()

	c := newRetryTestClient(t, srv.URL)
	_, err := handlers.DispatchRaw(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: "/test",
	})
	if err != nil {
		t.Fatalf("expected success after retry, got: %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Errorf("expected 3 attempts, got %d", got)
	}
}

func TestRetry_502_Retried(t *testing.T) {
	output.Format = "json"
	var attempts int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n == 1 {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(`bad gateway`))
			return
		}
		w.Header().Set(headerContentType, contentTypeJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1}`))
	}))
	defer srv.Close()

	c := newRetryTestClient(t, srv.URL)
	_, err := handlers.DispatchRaw(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: "/test",
	})
	if err != nil {
		t.Fatalf("expected success after retry, got: %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Errorf("expected 2 attempts, got %d", got)
	}
}

func TestRetry_504_Retried(t *testing.T) {
	output.Format = "json"
	var attempts int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n == 1 {
			w.WriteHeader(http.StatusGatewayTimeout)
			_, _ = w.Write([]byte(`gateway timeout`))
			return
		}
		w.Header().Set(headerContentType, contentTypeJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1}`))
	}))
	defer srv.Close()

	c := newRetryTestClient(t, srv.URL)
	_, err := handlers.DispatchRaw(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: "/test",
	})
	if err != nil {
		t.Fatalf("expected success after retry, got: %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Errorf("expected 2 attempts, got %d", got)
	}
}

func TestRetry_PermanentError_NoRetry(t *testing.T) {
	output.Format = "json"
	var attempts int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"message":"not found"}}`))
	}))
	defer srv.Close()

	c := newRetryTestClient(t, srv.URL)
	_, err := handlers.DispatchRaw(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: "/test",
	})
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Errorf("expected 1 attempt (no retry for 404), got %d", got)
	}
}

func TestRetry_200_NoRetry(t *testing.T) {
	output.Format = "json"
	var attempts int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.Header().Set(headerContentType, contentTypeJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1}`))
	}))
	defer srv.Close()

	c := newRetryTestClient(t, srv.URL)
	_, err := handlers.DispatchRaw(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: "/test",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Errorf("expected 1 attempt (no retry needed), got %d", got)
	}
}

func TestRetry_ExhaustedRetries_ReturnsError(t *testing.T) {
	output.Format = "json"
	var attempts int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":{"message":"service unavailable"}}`))
	}))
	defer srv.Close()

	c := newRetryTestClient(t, srv.URL)
	_, err := handlers.DispatchRaw(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: "/test",
	})
	if err == nil {
		t.Fatal("expected error after exhausted retries, got nil")
	}
	// Initial attempt + 3 retries = 4 total attempts
	if got := atomic.LoadInt32(&attempts); got != 4 {
		t.Errorf("expected 4 attempts (1 + 3 retries), got %d", got)
	}
}
