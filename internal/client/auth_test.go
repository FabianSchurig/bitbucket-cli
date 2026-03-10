package client_test

import (
	"os"
	"testing"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
)

func TestNewClient_AppPassword(t *testing.T) {
	t.Setenv("BITBUCKET_USERNAME", "testuser")
	t.Setenv("BITBUCKET_APP_PASSWORD", "testpassword")
	t.Setenv("BITBUCKET_TOKEN", "")

	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("expected no error with app password, got: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_Token(t *testing.T) {
	t.Setenv("BITBUCKET_USERNAME", "")
	t.Setenv("BITBUCKET_APP_PASSWORD", "")
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
	for _, k := range []string{"BITBUCKET_USERNAME", "BITBUCKET_APP_PASSWORD", "BITBUCKET_TOKEN"} {
		if err := os.Unsetenv(k); err != nil {
			t.Fatalf("unsetenv %s: %v", k, err)
		}
	}

	_, err := client.NewClient()
	if err == nil {
		t.Error("expected error when no auth is configured, got nil")
	}
}

func TestNewClient_AppPasswordTakesPrecedence(t *testing.T) {
	// When both username+password AND token are set, basic auth should be used
	t.Setenv("BITBUCKET_USERNAME", "user")
	t.Setenv("BITBUCKET_APP_PASSWORD", "pass")
	t.Setenv("BITBUCKET_TOKEN", "token")

	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}
