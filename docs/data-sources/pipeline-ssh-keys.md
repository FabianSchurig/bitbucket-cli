---
page_title: "bitbucket_pipeline_ssh_keys Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pipeline-ssh-keys via the Bitbucket Cloud API.
---

# bitbucket_pipeline_ssh_keys (Data Source)

Reads Bitbucket pipeline-ssh-keys via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-ssh-key-pair-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pipeline:bitbucket` |

## Example Usage

```hcl
data "bitbucket_pipeline_ssh_keys" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "pipeline_ssh_keys_response" {
  value = data.bitbucket_pipeline_ssh_keys.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
