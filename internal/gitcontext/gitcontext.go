// Package gitcontext infers Bitbucket workspace and repository slug from the
// current directory's git remote URL. It checks environment variables first,
// then falls back to parsing the git config for the "origin" remote.
package gitcontext

import (
	"bufio"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// InferDefaults returns workspace and repo slug values inferred from
// environment variables (BITBUCKET_WORKSPACE, BITBUCKET_REPO_SLUG) and
// the "origin" remote of the current directory's git repository.
//
// Precedence: environment variables > git remote URL.
// CLI flags (checked by the caller) take the highest precedence.
func InferDefaults() (workspace, repoSlug string) {
	workspace = os.Getenv("BITBUCKET_WORKSPACE")
	repoSlug = os.Getenv("BITBUCKET_REPO_SLUG")
	if workspace != "" && repoSlug != "" {
		return
	}

	gitWs, gitSlug := inferFromGitRemote()
	if workspace == "" {
		workspace = gitWs
	}
	if repoSlug == "" {
		repoSlug = gitSlug
	}
	return
}

// inferFromGitRemote finds the git directory, reads the remote "origin" URL,
// and parses the Bitbucket workspace and repo slug from it.
func inferFromGitRemote() (workspace, repoSlug string) {
	remoteURL := findRemoteURL("origin")
	if remoteURL == "" {
		return "", ""
	}
	return ParseBitbucketURL(remoteURL)
}

// findRemoteURL locates the .git directory and parses its config for the
// given remote's URL.
func findRemoteURL(remoteName string) string {
	gitDir := findGitDir()
	if gitDir == "" {
		return ""
	}
	configPath := filepath.Join(gitDir, "config")
	return parseGitConfigRemoteURL(configPath, remoteName)
}

// findGitDir walks up from the current working directory looking for a .git
// directory (or .git file for worktrees/submodules).
func findGitDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		gitPath := filepath.Join(dir, ".git")
		info, err := os.Stat(gitPath)
		if err == nil {
			if info.IsDir() {
				return gitPath
			}
			// .git may be a file (worktrees/submodules) with "gitdir: <path>"
			if resolved := resolveGitFile(gitPath); resolved != "" {
				return resolved
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// resolveGitFile reads a .git file (used by worktrees and submodules) and
// resolves the gitdir path it points to.
func resolveGitFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(string(data))
	if !strings.HasPrefix(line, "gitdir: ") {
		return ""
	}
	gitdir := strings.TrimPrefix(line, "gitdir: ")
	if !filepath.IsAbs(gitdir) {
		gitdir = filepath.Join(filepath.Dir(path), gitdir)
	}
	if info, err := os.Stat(gitdir); err == nil && info.IsDir() {
		return gitdir
	}
	return ""
}

// parseGitConfigRemoteURL reads a git config file and extracts the URL for
// the named remote. It does minimal INI-style parsing — enough for standard
// git config files.
func parseGitConfigRemoteURL(configPath, remoteName string) string {
	f, err := os.Open(configPath)
	if err != nil {
		return ""
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	inRemote := false
	target := `[remote "` + remoteName + `"]`

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == target {
			inRemote = true
			continue
		}
		if inRemote {
			if strings.HasPrefix(line, "[") {
				return ""
			}
			if strings.HasPrefix(line, "url = ") {
				return strings.TrimPrefix(line, "url = ")
			}
		}
	}
	return ""
}

// ParseBitbucketURL extracts the workspace and repository slug from a
// Bitbucket remote URL. It supports:
//   - SSH:   git@bitbucket.org:workspace/repo.git
//   - HTTPS: https://bitbucket.org/workspace/repo.git
//   - HTTPS with user: https://user@bitbucket.org/workspace/repo.git
//
// Returns empty strings if the URL is not a recognised Bitbucket format.
func ParseBitbucketURL(rawURL string) (workspace, repoSlug string) {
	// SCP-style SSH format: git@bitbucket.org:workspace/repo.git
	if !strings.Contains(rawURL, "://") {
		parts := strings.SplitN(rawURL, ":", 2)
		if len(parts) == 2 {
			host := parts[0]
			if at := strings.LastIndex(host, "@"); at >= 0 {
				host = host[at+1:]
			}
			if host == "bitbucket.org" {
				return extractFromPath(parts[1])
			}
		}
		return "", ""
	}

	// Also handle ssh://git@bitbucket.org/workspace/repo.git and HTTPS forms.
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", ""
	}

	host := u.Hostname()
	if host != "bitbucket.org" {
		return "", ""
	}
	return extractFromPath(u.Path)
}

// extractFromPath splits a "workspace/repo.git" path into workspace and slug.
func extractFromPath(path string) (workspace, repoSlug string) {
	path = strings.TrimSuffix(path, ".git")
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	parts := strings.SplitN(path, "/", 3)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", ""
	}
	return parts[0], parts[1]
}
