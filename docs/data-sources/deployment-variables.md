---
page_title: "bitbucket_deployment_variables Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket deployment-variables via the Bitbucket Cloud API.
---

# bitbucket_deployment_variables (Data Source)

Reads Bitbucket deployment-variables via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-deployments-config-environments-environment-uuid-variables-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pipeline:bitbucket` |

## Example Usage

```hcl
data "bitbucket_deployment_variables" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  environment_uuid = "env-uuid"
}

output "deployment_variables_response" {
  value = data.bitbucket_deployment_variables.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `environment_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `key` (String) The unique name of the variable.
- `secured` (String) If true, this variable will be treated as secured. The value will never be exposed in the logs or the REST API.
- `uuid` (String) The UUID identifying the variable.
- `value` (String) The value of the variable. If the variable is secured, this will be empty.
