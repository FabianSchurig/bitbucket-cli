---
page_title: "bitbucket_repo_runners Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket repo-runners via the Bitbucket Cloud API.
---

# bitbucket_repo_runners (Resource)

Manages Bitbucket repo-runners via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/pipelines-config/runners` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-runners-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines-config/runners/{runner_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-runners-runner-uuid-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/pipelines-config/runners/{runner_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-runners-runner-uuid-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/pipelines-config/runners/{runner_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-runners-runner-uuid-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines-config/runners` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-runners-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `write:runner:bitbucket`, `read:runner:bitbucket` |
| Read | `read:runner:bitbucket` |
| Update | `read:runner:bitbucket`, `write:runner:bitbucket` |
| Delete | `write:runner:bitbucket` |
| List | `read:runner:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_repo_runners" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  runner_uuid = "{runner-uuid}"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `runner_uuid` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
