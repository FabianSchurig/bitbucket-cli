---
page_title: "bitbucket_pipeline_oidc Resource - bitbucket"
subcategory: "Pipelines"
description: |-
  Manages Bitbucket pipeline-oidc via the Bitbucket Cloud API.
---

# bitbucket_pipeline_oidc (Resource)

Manages Bitbucket pipeline-oidc via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/pipelines-config/identity/oidc/.well-known/openid-configuration` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-workspaces-workspace-pipelines-config-identity-oidc-.well-known-openid-configuration-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | — |

## Example Usage

```hcl
resource "bitbucket_pipeline_oidc" "example" {
  workspace = "my-workspace"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
