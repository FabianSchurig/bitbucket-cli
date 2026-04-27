package client_test

import (
	"os"
	"testing"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
)

func TestNewClient_TokenOnly(t *testing.T) {
	t.Setenv("BITBUCKET_USERNAME", "")
	t.Setenv("BITBUCKET_TOKEN", "mytoken")

	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("expected no error with token, got: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_UsernameAndToken(t *testing.T) {
	t.Setenv("BITBUCKET_USERNAME", "testuser")
	t.Setenv("BITBUCKET_TOKEN", "testtoken")

	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("expected no error with username+token, got: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_NoAuth(t *testing.T) {
	for _, k := range []string{
		"BITBUCKET_USERNAME",
		"BITBUCKET_TOKEN",
		"BITBUCKET_CSRF_TOKEN",
		"BITBUCKET_CLOUD_SESSION_TOKEN",
	} {
		if err := os.Unsetenv(k); err != nil {
			t.Fatalf("unsetenv %s: %v", k, err)
		}
	}

	_, err := client.NewClient()
	if err == nil {
		t.Error("expected error when no auth is configured, got nil")
	}
}

// TestNewClient_InternalAPITokensOnly verifies that the client can be
// constructed with only the cookie-based auth tokens used by Bitbucket's
// internal API (csrftoken + cloud.session.token), without BITBUCKET_TOKEN.
//
// Internal endpoints (https://bitbucket.org/!api/internal/...) do not accept
// HTTP Basic Auth — they require the same browser-style cookies + CSRF header
// the Bitbucket UI sends. This is the only way to talk to them.
func TestNewClient_InternalAPITokensOnly(t *testing.T) {
	t.Setenv("BITBUCKET_USERNAME", "")
	t.Setenv("BITBUCKET_TOKEN", "")
	t.Setenv("BITBUCKET_CSRF_TOKEN", "csrf-abc")
	t.Setenv("BITBUCKET_CLOUD_SESSION_TOKEN", "session-xyz")

	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("expected no error with internal-API tokens, got: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.CSRFToken != "csrf-abc" {
		t.Errorf("CSRFToken = %q, want %q", c.CSRFToken, "csrf-abc")
	}
	if c.CloudSessionToken != "session-xyz" {
		t.Errorf("CloudSessionToken = %q, want %q", c.CloudSessionToken, "session-xyz")
	}
}

// TestNewClientWithConfig_StoresAllCredentials verifies the explicit
// configuration constructor accepts and stores the internal-API cookies in
// addition to username/token, so that callers like the Terraform provider can
// pass them through without touching environment variables.
func TestNewClientWithConfig_StoresAllCredentials(t *testing.T) {
	c, err := client.NewClientWithConfig("user", "tok", "", "csrf", "sess")
	if err != nil {
		t.Fatalf("NewClientWithConfig: %v", err)
	}
	if c.Username != "user" || c.Token != "tok" {
		t.Errorf("Username/Token = %q/%q, want user/tok", c.Username, c.Token)
	}
	if c.CSRFToken != "csrf" || c.CloudSessionToken != "sess" {
		t.Errorf("CSRFToken/CloudSessionToken = %q/%q, want csrf/sess",
			c.CSRFToken, c.CloudSessionToken)
	}
}
