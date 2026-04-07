package handlers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInferRepoContext_BothEmpty(t *testing.T) {
	setupFakeGitRepo(t, "git@bitbucket.org:myws/myrepo.git")
	t.Setenv("BITBUCKET_WORKSPACE", "")
	t.Setenv("BITBUCKET_REPO_SLUG", "")

	pathParams := map[string]string{
		"workspace": "",
		"repo_slug": "",
	}
	InferRepoContext(pathParams)

	if pathParams["workspace"] != "myws" {
		t.Errorf("workspace = %q, want %q", pathParams["workspace"], "myws")
	}
	if pathParams["repo_slug"] != "myrepo" {
		t.Errorf("repo_slug = %q, want %q", pathParams["repo_slug"], "myrepo")
	}
}

func TestInferRepoContext_AlreadySet(t *testing.T) {
	setupFakeGitRepo(t, "git@bitbucket.org:myws/myrepo.git")
	t.Setenv("BITBUCKET_WORKSPACE", "")
	t.Setenv("BITBUCKET_REPO_SLUG", "")

	pathParams := map[string]string{
		"workspace": "explicit",
		"repo_slug": "given",
	}
	InferRepoContext(pathParams)

	if pathParams["workspace"] != "explicit" {
		t.Errorf("workspace = %q, want %q", pathParams["workspace"], "explicit")
	}
	if pathParams["repo_slug"] != "given" {
		t.Errorf("repo_slug = %q, want %q", pathParams["repo_slug"], "given")
	}
}

func TestInferRepoContext_Partial(t *testing.T) {
	setupFakeGitRepo(t, "git@bitbucket.org:myws/myrepo.git")
	t.Setenv("BITBUCKET_WORKSPACE", "")
	t.Setenv("BITBUCKET_REPO_SLUG", "")

	pathParams := map[string]string{
		"workspace": "explicit",
		"repo_slug": "",
	}
	InferRepoContext(pathParams)

	if pathParams["workspace"] != "explicit" {
		t.Errorf("workspace = %q, want %q", pathParams["workspace"], "explicit")
	}
	if pathParams["repo_slug"] != "myrepo" {
		t.Errorf("repo_slug = %q, want %q", pathParams["repo_slug"], "myrepo")
	}
}

func TestInferRepoContext_NoRelevantKeys(t *testing.T) {
	// Map without workspace/repo_slug — should be a no-op.
	pathParams := map[string]string{
		"commit": "abc123",
	}
	InferRepoContext(pathParams)

	if pathParams["commit"] != "abc123" {
		t.Errorf("commit = %q, want %q", pathParams["commit"], "abc123")
	}
	if _, ok := pathParams["workspace"]; ok {
		t.Error("workspace should not have been added to map")
	}
}

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
