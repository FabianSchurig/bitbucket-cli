---
page_title: "bitbucket_workspace_runners Data Source - bitbucket"
subcategory: ""
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
  runner_uuid = "{runner-uuid}"
}

output "workspace_runners_response" {
  value = data.bitbucket_workspace_runners.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `runner_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the runner was created.
- `labels` (String) Labels assigned to the runner for identification and routing. (JSON array)
- `name` (String) The name of the runner.
- `oauth_client_audience` (String) The intended audience for the OAuth token.
- `oauth_client_id` (String) The OAuth client ID.
- `oauth_client_secret` (String) The OAuth client secret. This is an optional element that is only provided once.
- `oauth_client_token_endpoint` (String) The OAuth token endpoint URL.
- `state_cordoned` (String) Whether the runner is cordoned (prevented from accepting new steps).
- `state_status` (String) The current status of the runner. [UNREGISTERED, ONLINE, OFFLINE, DISABLED, ENABLED, UNHEALTHY]
- `state_updated_on` (String) The timestamp when the runner state was last updated.
- `state_version_current` (String) The current recommended version of the runner.
- `state_version_version` (String) The currently installed version of the runner.
- `updated_on` (String) The timestamp when the runner was last updated.
- `uuid` (String) The UUID identifying the runner.
