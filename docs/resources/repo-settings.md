---
page_title: "bitbucket_repo_settings Resource - bitbucket"
subcategory: "Repositories"
description: |-
  Manages Bitbucket repo-settings via the Bitbucket Cloud API.
---

# bitbucket_repo_settings (Resource)

Manages Bitbucket repo-settings via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `PUT` | `/repositories/{workspace}/{repo_slug}/override-settings` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-override-settings-put) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/override-settings` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-override-settings-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/override-settings` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-override-settings-put) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `admin:repository:bitbucket` |
| Read | `admin:repository:bitbucket` |
| Update | `admin:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_repo_settings" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `type` (String) type

## Import

Existing resources can be imported into Terraform state. The import ID is the
slash-separated list of path parameter values in URL order: `workspace/repo_slug`.

Using an `import` block (Terraform 1.5+):

```hcl
import {
  to = bitbucket_repo_settings.example
  id = "my-workspace/my-repo"
}
```

Or with the CLI:

```shell
terraform import bitbucket_repo_settings.example "my-workspace/my-repo"
```
