---
page_title: "bitbucket_commit_file Resource - bitbucket"
subcategory: "Repositories"
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
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `commit` (String) Path parameter (auto-populated from API response).
- `path` (String) Path parameter (auto-populated from API response).

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `type` (String) type

## Import

Existing resources can be imported into Terraform state. The import ID is the
slash-separated list of path parameter values in URL order: `workspace/repo_slug/commit/path`.

Using an `import` block (Terraform 1.5+):

```hcl
import {
  to = bitbucket_commit_file.example
  id = "my-workspace/my-repo/abc123def/README.md"
}
```

Or with the CLI:

```shell
terraform import bitbucket_commit_file.example "my-workspace/my-repo/abc123def/README.md"
```
