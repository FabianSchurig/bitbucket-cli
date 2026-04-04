---
page_title: "bitbucket_pipeline_oidc_keys Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket pipeline-oidc-keys via the Bitbucket Cloud API.
---

# bitbucket_pipeline_oidc_keys (Resource)

Manages Bitbucket pipeline-oidc-keys via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/pipelines-config/identity/oidc/keys.json` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-workspaces-workspace-pipelines-config-identity-oidc-keys.json-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | — |

## Example Usage

```hcl
resource "bitbucket_pipeline_oidc_keys" "example" {
  workspace = "my-workspace"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
