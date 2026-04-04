---
page_title: "bitbucket_branching_model Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket branching-model via the Bitbucket Cloud API.
---

# bitbucket_branching_model (Resource)

Manages Bitbucket branching-model via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/branching-model` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branching-model/#api-repositories-workspace-repo-slug-branching-model-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/branching-model/settings` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branching-model/#api-repositories-workspace-repo-slug-branching-model-settings-put) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| Update | `admin:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_branching_model" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
