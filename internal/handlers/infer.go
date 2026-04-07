package handlers

import "github.com/FabianSchurig/bitbucket-cli/internal/gitcontext"

// InferRepoContext fills empty "workspace" and "repo_slug" entries in
// pathParams from environment variables (BITBUCKET_WORKSPACE,
// BITBUCKET_REPO_SLUG) and the current directory's git remote URL.
//
// Only keys already present in the map are considered; commands that
// don't use workspace or repo_slug are unaffected.
func InferRepoContext(pathParams map[string]string) {
	_, hasWs := pathParams["workspace"]
	_, hasSlug := pathParams["repo_slug"]
	if !hasWs && !hasSlug {
		return
	}

	needWs := hasWs && pathParams["workspace"] == ""
	needSlug := hasSlug && pathParams["repo_slug"] == ""
	if !needWs && !needSlug {
		return
	}

	ws, slug := gitcontext.InferDefaults()
	if needWs {
		pathParams["workspace"] = ws
	}
	if needSlug {
		pathParams["repo_slug"] = slug
	}
}
