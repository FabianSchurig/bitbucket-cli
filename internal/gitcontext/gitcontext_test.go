package gitcontext

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseBitbucketURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantWS   string
		wantSlug string
	}{
		{
			name:     "SSH standard",
			url:      "git@bitbucket.org:myteam/myrepo.git",
			wantWS:   "myteam",
			wantSlug: "myrepo",
		},
		{
			name:     "SSH without .git",
			url:      "git@bitbucket.org:myteam/myrepo",
			wantWS:   "myteam",
			wantSlug: "myrepo",
		},
		{
			name:     "HTTPS standard",
			url:      "https://bitbucket.org/myteam/myrepo.git",
			wantWS:   "myteam",
			wantSlug: "myrepo",
		},
		{
			name:     "HTTPS without .git",
			url:      "https://bitbucket.org/myteam/myrepo",
			wantWS:   "myteam",
			wantSlug: "myrepo",
		},
		{
			name:     "HTTPS with user",
			url:      "https://user@bitbucket.org/myteam/myrepo.git",
			wantWS:   "myteam",
			wantSlug: "myrepo",
		},
		{
			name:     "SSH with ssh:// scheme",
			url:      "ssh://git@bitbucket.org/myteam/myrepo.git",
			wantWS:   "myteam",
			wantSlug: "myrepo",
		},
		{
			name:     "not bitbucket",
			url:      "git@github.com:owner/repo.git",
			wantWS:   "",
			wantSlug: "",
		},
		{
			name:     "HTTPS not bitbucket",
			url:      "https://github.com/owner/repo.git",
			wantWS:   "",
			wantSlug: "",
		},
		{
			name:     "SCP spoofed host suffix",
			url:      "git@evilbitbucket.org:ws/repo.git",
			wantWS:   "",
			wantSlug: "",
		},
		{
			name:     "HTTPS spoofed host suffix",
			url:      "https://evilbitbucket.org/ws/repo.git",
			wantWS:   "",
			wantSlug: "",
		},
		{
			name:     "empty string",
			url:      "",
			wantWS:   "",
			wantSlug: "",
		},
		{
			name:     "only workspace",
			url:      "git@bitbucket.org:myteam",
			wantWS:   "",
			wantSlug: "",
		},
		{
			name:     "trailing slash",
			url:      "https://bitbucket.org/myteam/myrepo/",
			wantWS:   "myteam",
			wantSlug: "myrepo",
		},
		{
			name:     "extra path segments",
			url:      "https://bitbucket.org/myteam/myrepo/src/main",
			wantWS:   "myteam",
			wantSlug: "myrepo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws, slug := ParseBitbucketURL(tt.url)
			if ws != tt.wantWS || slug != tt.wantSlug {
				t.Errorf("ParseBitbucketURL(%q) = (%q, %q), want (%q, %q)",
					tt.url, ws, slug, tt.wantWS, tt.wantSlug)
			}
		})
	}
}

func TestParseGitConfigRemoteURL(t *testing.T) {
	configContent := `[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
[remote "origin"]
	url = git@bitbucket.org:myteam/myrepo.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[remote "upstream"]
	url = https://bitbucket.org/otherteam/otherrepo.git
	fetch = +refs/heads/*:refs/remotes/upstream/*
[branch "main"]
	remote = origin
	merge = refs/heads/main
`
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config")
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		remote  string
		wantURL string
	}{
		{"origin", "git@bitbucket.org:myteam/myrepo.git"},
		{"upstream", "https://bitbucket.org/otherteam/otherrepo.git"},
		{"nonexistent", ""},
	}

	for _, tt := range tests {
		t.Run(tt.remote, func(t *testing.T) {
			got := parseGitConfigRemoteURL(configPath, tt.remote)
			if got != tt.wantURL {
				t.Errorf("parseGitConfigRemoteURL(%q, %q) = %q, want %q",
					configPath, tt.remote, got, tt.wantURL)
			}
		})
	}
}

func TestParseGitConfigRemoteURL_MissingFile(t *testing.T) {
	got := parseGitConfigRemoteURL("/nonexistent/path/config", "origin")
	if got != "" {
		t.Errorf("expected empty string for missing file, got %q", got)
	}
}

func TestFindGitDir(t *testing.T) {
	// Create a temporary directory structure with a .git dir
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	if err := os.Mkdir(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(root, "sub", "deep")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Change to the subdirectory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	if err := os.Chdir(subDir); err != nil {
		t.Fatal(err)
	}

	got := findGitDir()
	if got != gitDir {
		t.Errorf("findGitDir() = %q, want %q", got, gitDir)
	}
}

func TestFindGitDir_WorktreeFile(t *testing.T) {
	// Simulate a .git file pointing to another directory (worktree)
	root := t.TempDir()
	realGitDir := filepath.Join(root, "real-git-dir")
	if err := os.MkdirAll(realGitDir, 0o755); err != nil {
		t.Fatal(err)
	}

	worktree := filepath.Join(root, "worktree")
	if err := os.MkdirAll(worktree, 0o755); err != nil {
		t.Fatal(err)
	}
	gitFile := filepath.Join(worktree, ".git")
	if err := os.WriteFile(gitFile, []byte("gitdir: "+realGitDir+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	if err := os.Chdir(worktree); err != nil {
		t.Fatal(err)
	}

	got := findGitDir()
	if got != realGitDir {
		t.Errorf("findGitDir() = %q, want %q", got, realGitDir)
	}
}

func TestInferDefaults_EnvVarsOnly(t *testing.T) {
	t.Setenv("BITBUCKET_WORKSPACE", "envws")
	t.Setenv("BITBUCKET_REPO_SLUG", "envrepo")

	ws, slug := InferDefaults()
	if ws != "envws" || slug != "envrepo" {
		t.Errorf("InferDefaults() = (%q, %q), want (envws, envrepo)", ws, slug)
	}
}

func TestInferDefaults_EnvVarsPartial(t *testing.T) {
	// Set only workspace via env; repo-slug should come from git if available
	t.Setenv("BITBUCKET_WORKSPACE", "envws")
	t.Setenv("BITBUCKET_REPO_SLUG", "")

	// No git repo in temp dir, so repoSlug stays empty
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}

	ws, slug := InferDefaults()
	if ws != "envws" {
		t.Errorf("workspace = %q, want envws", ws)
	}
	if slug != "" {
		t.Errorf("repoSlug = %q, want empty (no git repo)", slug)
	}
}

func TestInferDefaults_GitRemote(t *testing.T) {
	t.Setenv("BITBUCKET_WORKSPACE", "")
	t.Setenv("BITBUCKET_REPO_SLUG", "")

	// Create a fake git repo with a config
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	if err := os.Mkdir(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	config := `[remote "origin"]
	url = git@bitbucket.org:myteam/myrepo.git
`
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

	ws, slug := InferDefaults()
	if ws != "myteam" || slug != "myrepo" {
		t.Errorf("InferDefaults() = (%q, %q), want (myteam, myrepo)", ws, slug)
	}
}

func TestExtractFromPath(t *testing.T) {
	tests := []struct {
		path     string
		wantWS   string
		wantSlug string
	}{
		{"myteam/myrepo.git", "myteam", "myrepo"},
		{"/myteam/myrepo.git", "myteam", "myrepo"},
		{"myteam/myrepo", "myteam", "myrepo"},
		{"/myteam/myrepo/", "myteam", "myrepo"},
		{"myteam", "", ""},
		{"", "", ""},
		{"/", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			ws, slug := extractFromPath(tt.path)
			if ws != tt.wantWS || slug != tt.wantSlug {
				t.Errorf("extractFromPath(%q) = (%q, %q), want (%q, %q)",
					tt.path, ws, slug, tt.wantWS, tt.wantSlug)
			}
		})
	}
}
