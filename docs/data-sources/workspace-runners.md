---
page_title: "bitbucket_workspace_runners Data Source - bitbucket"
subcategory: "Pipelines"
description: |-
  Reads Bitbucket workspace-runners via the Bitbucket Cloud API.
---

# bitbucket_workspace_runners (Data Source)

Reads Bitbucket workspace-runners via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/pipelines-config/runners/{runner_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-workspaces-workspace-pipelines-config-runners-runner-uuid-get) |
| List | `GET` | `/workspaces/{workspace}/pipelines-config/runners` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-workspaces-workspace-pipelines-config-runners-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:runner:bitbucket` |
| List | `read:runner:bitbucket` |

## Example Usage

```hcl
data "bitbucket_workspace_runners" "example" {
  workspace = "my-workspace"
}

output "workspace_runners_response" {
  value = data.bitbucket_workspace_runners.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Optional
- `runner_uuid` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the runner was created.
- `labels` (List of String) Labels assigned to the runner for identification and routing.
- `name` (String) The name of the runner.
- `oauth_client` (Object) oauth_client
  Nested schema:
  - `id` (String) The OAuth client ID.
  - `secret` (String) The OAuth client secret. This is an optional element that is only provided once.
  - `token_endpoint` (String) The OAuth token endpoint URL.
  - `audience` (String) The intended audience for the OAuth token.

- `state` (Object) state
  Nested schema:
  - `updated_on` (String) The timestamp when the runner state was last updated.
  - `cordoned` (String) Whether the runner is cordoned (prevented from accepting new steps).
  - `status` (String) The current status of the runner. [UNREGISTERED, ONLINE, OFFLINE, DISABLED, ENABLED, UNHEALTHY]

- `updated_on` (String) The timestamp when the runner was last updated.
- `uuid` (String) The UUID identifying the runner.
