package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

// setupFakeGitRepo creates a temporary directory containing a minimal .git/config
// that points at the given Bitbucket remote URL. It changes the working directory
// to that repository root and restores the original directory on cleanup.
func setupFakeGitRepo(t *testing.T, remoteURL string) {
	t.Helper()
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	if err := os.Mkdir(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	config := "[remote \"origin\"]\n\turl = " + remoteURL + "\n"
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
}

// setupInferenceEnv sets auth and base-URL env vars for the test server and
// clears the workspace/repo-slug env vars so inference can be tested.
func setupInferenceEnv(t *testing.T, baseURL string) {
	t.Helper()
	t.Setenv("BITBUCKET_TOKEN", "test-token")
	t.Setenv("BITBUCKET_USERNAME", "testuser")
	t.Setenv("BITBUCKET_BASE_URL", baseURL)
	t.Setenv("BITBUCKET_WORKSPACE", "")
	t.Setenv("BITBUCKET_REPO_SLUG", "")
}

// TestInference_GitRemote_SSH verifies that workspace and repo-slug are
// inferred from an SSH git remote when the flags are omitted.
func TestInference_GitRemote_SSH(t *testing.T) {
	output.Format = "json"
	setupFakeGitRepo(t, "git@bitbucket.org:inferred-ws/inferred-repo.git")

	var receivedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"hash": "abc123"})
	}))
	defer srv.Close()

	setupInferenceEnv(t, srv.URL)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"commits", "get-a-commit", "--commit", "abc123"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// Verify the request path contains the inferred workspace and repo slug.
	if !strings.Contains(receivedPath, "/repositories/inferred-ws/inferred-repo/") {
		t.Errorf("expected inferred workspace/repo in path, got %s", receivedPath)
	}
}

// TestInference_GitRemote_HTTPS verifies inference from an HTTPS remote.
func TestInference_GitRemote_HTTPS(t *testing.T) {
	output.Format = "json"
	setupFakeGitRepo(t, "https://bitbucket.org/https-ws/https-repo.git")

	var receivedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"hash": "abc123"})
	}))
	defer srv.Close()

	setupInferenceEnv(t, srv.URL)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"commits", "get-a-commit", "--commit", "abc123"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	if !strings.Contains(receivedPath, "/repositories/https-ws/https-repo/") {
		t.Errorf("expected HTTPS-inferred workspace/repo in path, got %s", receivedPath)
	}
}

// TestInference_FlagOverridesGit verifies that explicit --workspace and
// --repo-slug flags take precedence over the git remote.
func TestInference_FlagOverridesGit(t *testing.T) {
	output.Format = "json"
	// Git remote points to inferred-ws/inferred-repo
	setupFakeGitRepo(t, "git@bitbucket.org:inferred-ws/inferred-repo.git")

	var receivedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"hash": "abc123"})
	}))
	defer srv.Close()

	setupInferenceEnv(t, srv.URL)

	cmd := newRootCmd()
	cmd.SetArgs([]string{
		"commits", "get-a-commit",
		"--commit", "abc123",
		"--workspace", "explicit-ws",
		"--repo-slug", "explicit-repo",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// The explicit flags should override the git remote.
	if !strings.Contains(receivedPath, "/repositories/explicit-ws/explicit-repo/") {
		t.Errorf("expected explicit flags to override git inference, got path %s", receivedPath)
	}
}

// TestInference_EnvVarOverridesGit verifies that BITBUCKET_WORKSPACE and
// BITBUCKET_REPO_SLUG environment variables take precedence over git remote.
func TestInference_EnvVarOverridesGit(t *testing.T) {
	output.Format = "json"
	// Git remote points to inferred-ws/inferred-repo
	setupFakeGitRepo(t, "git@bitbucket.org:inferred-ws/inferred-repo.git")

	var receivedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"hash": "abc123"})
	}))
	defer srv.Close()

	setupInferenceEnv(t, srv.URL)
	// Override with env vars — these should win over git remote.
	t.Setenv("BITBUCKET_WORKSPACE", "env-ws")
	t.Setenv("BITBUCKET_REPO_SLUG", "env-repo")

	cmd := newRootCmd()
	cmd.SetArgs([]string{"commits", "get-a-commit", "--commit", "abc123"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	if !strings.Contains(receivedPath, "/repositories/env-ws/env-repo/") {
		t.Errorf("expected env vars to override git inference, got path %s", receivedPath)
	}
}

// TestInference_FlagOverridesEnvVar verifies that explicit flags take
// precedence over environment variables.
func TestInference_FlagOverridesEnvVar(t *testing.T) {
	output.Format = "json"
	setupFakeGitRepo(t, "git@bitbucket.org:inferred-ws/inferred-repo.git")

	var receivedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"hash": "abc123"})
	}))
	defer srv.Close()

	setupInferenceEnv(t, srv.URL)
	t.Setenv("BITBUCKET_WORKSPACE", "env-ws")
	t.Setenv("BITBUCKET_REPO_SLUG", "env-repo")

	cmd := newRootCmd()
	cmd.SetArgs([]string{
		"commits", "get-a-commit",
		"--commit", "abc123",
		"--workspace", "flag-ws",
		"--repo-slug", "flag-repo",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	if !strings.Contains(receivedPath, "/repositories/flag-ws/flag-repo/") {
		t.Errorf("expected flags to override env vars, got path %s", receivedPath)
	}
}

// TestInference_PartialOverride verifies that one flag can be inferred while
// the other is provided explicitly.
func TestInference_PartialOverride(t *testing.T) {
	output.Format = "json"
	setupFakeGitRepo(t, "git@bitbucket.org:inferred-ws/inferred-repo.git")

	var receivedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"hash": "abc123"})
	}))
	defer srv.Close()

	setupInferenceEnv(t, srv.URL)

	// Only provide --workspace explicitly; repo-slug should be inferred.
	cmd := newRootCmd()
	cmd.SetArgs([]string{
		"commits", "get-a-commit",
		"--commit", "abc123",
		"--workspace", "explicit-ws",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	if !strings.Contains(receivedPath, "/repositories/explicit-ws/inferred-repo/") {
		t.Errorf("expected mixed explicit+inferred path, got %s", receivedPath)
	}
}

// TestInference_NoGitRepo_StillRequiresFlags verifies that when not inside a
// git repository and no env vars are set, the required flag error is returned.
func TestInference_NoGitRepo_StillRequiresFlags(t *testing.T) {
	output.Format = "json"

	// Use a temporary directory with no .git
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}

	t.Setenv("BITBUCKET_TOKEN", "test-token")
	t.Setenv("BITBUCKET_USERNAME", "testuser")
	t.Setenv("BITBUCKET_WORKSPACE", "")
	t.Setenv("BITBUCKET_REPO_SLUG", "")

	cmd := newRootCmd()
	cmd.SetArgs([]string{"commits", "get-a-commit", "--commit", "abc123"})
	err = cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing workspace/repo-slug when not in git repo")
	}
	if !strings.Contains(err.Error(), "is required") {
		t.Errorf("expected 'is required' error, got: %v", err)
	}
}

// TestInference_NonBitbucketRemote_StillRequiresFlags verifies that a
// non-Bitbucket remote (e.g. GitHub) does not produce inference.
func TestInference_NonBitbucketRemote_StillRequiresFlags(t *testing.T) {
	output.Format = "json"
	setupFakeGitRepo(t, "git@github.com:owner/repo.git")

	t.Setenv("BITBUCKET_TOKEN", "test-token")
	t.Setenv("BITBUCKET_USERNAME", "testuser")
	t.Setenv("BITBUCKET_WORKSPACE", "")
	t.Setenv("BITBUCKET_REPO_SLUG", "")

	cmd := newRootCmd()
	cmd.SetArgs([]string{"commits", "get-a-commit", "--commit", "abc123"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for non-Bitbucket remote")
	}
	if !strings.Contains(err.Error(), "is required") {
		t.Errorf("expected 'is required' error, got: %v", err)
	}
}
