---
page_title: "bitbucket_workspace_runners Resource - bitbucket"
subcategory: "Pipelines"
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
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Optional
- `runner_uuid` (String) Path parameter (auto-populated from API response).

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the runner was created.
- `labels` (List of String) Labels assigned to the runner for identification and routing.
- `name` (String) The name of the runner.
- `oauth_client` (Object) oauth_client
  Nested schema:
  - `audience` (String) The intended audience for the OAuth token.
  - `id` (String) The OAuth client ID.
  - `secret` (String) The OAuth client secret. This is an optional element that is only provided once.
  - `token_endpoint` (String) The OAuth token endpoint URL.

- `state` (Object) state
  Nested schema:
  - `cordoned` (String) Whether the runner is cordoned (prevented from accepting new steps).
  - `status` (String) The current status of the runner. [UNREGISTERED, ONLINE, OFFLINE, DISABLED, ENABLED, UNHEALTHY]
  - `updated_on` (String) The timestamp when the runner state was last updated.

- `updated_on` (String) The timestamp when the runner was last updated.
- `uuid` (String) The UUID identifying the runner.
