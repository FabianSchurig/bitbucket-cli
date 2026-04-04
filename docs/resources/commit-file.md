---
page_title: "bitbucket_commit_file Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket commit-file via the Bitbucket Cloud API.
---

# bitbucket_commit_file (Resource)

Manages Bitbucket commit-file via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/src` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-src-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/src/{commit}/{path}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-src-commit-path-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `write:repository:bitbucket` |
| Read | `read:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_commit_file" "example" {
  commit = "abc123def"
  path = "README.md"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `commit` (String) Path parameter.
- `path` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `path` (String) The path in the repository
- `type` (String) type
