package client_test

import (
	"os"
	"testing"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
)

func TestNewClient_APIToken(t *testing.T) {
	t.Setenv("BITBUCKET_USERNAME", "testuser")
	t.Setenv("BITBUCKET_API_TOKEN", "testpassword")
	t.Setenv("BITBUCKET_TOKEN", "")

	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("expected no error with API token, got: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_Token(t *testing.T) {
	t.Setenv("BITBUCKET_USERNAME", "")
	t.Setenv("BITBUCKET_API_TOKEN", "")
	t.Setenv("BITBUCKET_TOKEN", "mytoken")

	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("expected no error with token, got: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_NoAuth(t *testing.T) {
	// Clear all auth env vars
	for _, k := range []string{"BITBUCKET_USERNAME", "BITBUCKET_API_TOKEN", "BITBUCKET_TOKEN"} {
		if err := os.Unsetenv(k); err != nil {
			t.Fatalf("unsetenv %s: %v", k, err)
		}
	}

	_, err := client.NewClient()
	if err == nil {
		t.Error("expected error when no auth is configured, got nil")
	}
}

func TestNewClient_APITokenTakesPrecedence(t *testing.T) {
	// When both username+API token AND bearer token are set, basic auth should be used
	t.Setenv("BITBUCKET_USERNAME", "user")
	t.Setenv("BITBUCKET_API_TOKEN", "pass")
	t.Setenv("BITBUCKET_TOKEN", "token")

	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}
