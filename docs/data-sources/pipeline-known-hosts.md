---
page_title: "bitbucket_pipeline_known_hosts Data Source - bitbucket"
subcategory: "Pipelines"
description: |-
  Reads Bitbucket pipeline-known-hosts via the Bitbucket Cloud API.
---

# bitbucket_pipeline_known_hosts (Data Source)

Reads Bitbucket pipeline-known-hosts via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-ssh-known-hosts-known-host-uuid-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-ssh-known-hosts-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pipeline:bitbucket` |
| List | `read:pipeline:bitbucket` |

## Example Usage

```hcl
data "bitbucket_pipeline_known_hosts" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "pipeline_known_hosts_response" {
  value = data.bitbucket_pipeline_known_hosts.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `known_host_uuid` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `hostname` (String) The hostname of the known host.
- `public_key` (Object) public_key
  Nested schema:
  - `key` (String) The base64 encoded public key.
  - `key_type` (String) The type of the public key.
  - `md5_fingerprint` (String) The MD5 fingerprint of the public key.
  - `sha256_fingerprint` (String) The SHA-256 fingerprint of the public key.

- `uuid` (String) The UUID identifying the known host.
