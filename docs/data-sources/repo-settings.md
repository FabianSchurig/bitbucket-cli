---
page_title: "bitbucket_repo_settings Data Source - bitbucket"
subcategory: "Repositories"
description: |-
  Reads Bitbucket repo-settings via the Bitbucket Cloud API.
---

# bitbucket_repo_settings (Data Source)

Reads Bitbucket repo-settings via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/override-settings` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-override-settings-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `admin:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_repo_settings" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "repo_settings_response" {
  value = data.bitbucket_repo_settings.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `type` (String) type
