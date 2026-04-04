---
page_title: "bitbucket_pipeline_known_hosts Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket pipeline-known-hosts via the Bitbucket Cloud API.
---

# bitbucket_pipeline_known_hosts (Resource)

Manages Bitbucket pipeline-known-hosts via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-ssh-known-hosts-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-ssh-known-hosts-known-host-uuid-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-ssh-known-hosts-known-host-uuid-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-ssh-known-hosts-known-host-uuid-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-ssh-known-hosts-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `admin:pipeline:bitbucket` |
| Read | `read:pipeline:bitbucket` |
| Update | `admin:pipeline:bitbucket` |
| Delete | `admin:pipeline:bitbucket` |
| List | `read:pipeline:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_pipeline_known_hosts" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `known_host_uuid` (String) Path parameter (auto-populated from API response).
- `hostname` (String) The hostname of the known host. (also computed from API response)
- `public_key_key` (String) The base64 encoded public key. (also computed from API response)
- `public_key_key_type` (String) The type of the public key. (also computed from API response)
- `public_key_md5_fingerprint` (String) The MD5 fingerprint of the public key. (also computed from API response)
- `public_key_sha256_fingerprint` (String) The SHA-256 fingerprint of the public key. (also computed from API response)
- `uuid` (String) The UUID identifying the known host. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
