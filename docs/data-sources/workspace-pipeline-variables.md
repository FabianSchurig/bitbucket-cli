---
page_title: "bitbucket_workspace_pipeline_variables Data Source - bitbucket"
subcategory: "Pipelines"
description: |-
  Reads Bitbucket workspace-pipeline-variables via the Bitbucket Cloud API.
---

# bitbucket_workspace_pipeline_variables (Data Source)

Reads Bitbucket workspace-pipeline-variables via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/pipelines-config/variables/{variable_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-workspaces-workspace-pipelines-config-variables-variable-uuid-get) |
| List | `GET` | `/workspaces/{workspace}/pipelines-config/variables` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-workspaces-workspace-pipelines-config-variables-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pipeline:bitbucket` |
| List | `read:pipeline:bitbucket` |

## Example Usage

```hcl
data "bitbucket_workspace_pipeline_variables" "example" {
  workspace = "my-workspace"
}

output "workspace_pipeline_variables_response" {
  value = data.bitbucket_workspace_pipeline_variables.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Optional
- `variable_uuid` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `key` (String) The unique name of the variable.
- `secured` (String) If true, this variable will be treated as secured. The value will never be exposed in the logs or the REST API.
- `uuid` (String) The UUID identifying the variable.
- `value` (String) The value of the variable. If the variable is secured, this will be empty.
