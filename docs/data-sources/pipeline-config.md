---
page_title: "bitbucket_pipeline_config Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pipeline-config via the Bitbucket Cloud API.
---

# bitbucket_pipeline_config (Data Source)

Reads Bitbucket pipeline-config via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `admin:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_pipeline_config" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "pipeline_config_response" {
  value = data.bitbucket_pipeline_config.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `enabled` (String) Whether Pipelines is enabled for the repository.
