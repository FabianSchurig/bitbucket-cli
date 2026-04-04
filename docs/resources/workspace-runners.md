---
page_title: "bitbucket_workspace_runners Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket workspace-runners via the Bitbucket Cloud API.
---

# bitbucket_workspace_runners (Resource)

Manages Bitbucket workspace-runners via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/workspaces/{workspace}/pipelines-config/runners` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-workspaces-workspace-pipelines-config-runners-post) |
| Read | `GET` | `/workspaces/{workspace}/pipelines-config/runners/{runner_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-workspaces-workspace-pipelines-config-runners-runner-uuid-get) |
| Update | `PUT` | `/workspaces/{workspace}/pipelines-config/runners/{runner_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-workspaces-workspace-pipelines-config-runners-runner-uuid-put) |
| Delete | `DELETE` | `/workspaces/{workspace}/pipelines-config/runners/{runner_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-workspaces-workspace-pipelines-config-runners-runner-uuid-delete) |
| List | `GET` | `/workspaces/{workspace}/pipelines-config/runners` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-workspaces-workspace-pipelines-config-runners-get) |

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
resource "bitbucket_workspace_runners" "example" {
  workspace = "my-workspace"
  runner_uuid = "{runner-uuid}"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `runner_uuid` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
