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
