---
page_title: "bitbucket_pipeline_ssh_keys Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket pipeline-ssh-keys via the Bitbucket Cloud API.
---

# bitbucket_pipeline_ssh_keys (Resource)

Manages Bitbucket pipeline-ssh-keys via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-ssh-key-pair-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-ssh-key-pair-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-ssh-key-pair-delete) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pipeline:bitbucket` |
| Update | `admin:pipeline:bitbucket` |
| Delete | `admin:pipeline:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_pipeline_ssh_keys" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
