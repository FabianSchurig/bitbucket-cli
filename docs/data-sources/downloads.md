---
page_title: "bitbucket_downloads Data Source - bitbucket"
subcategory: "Downloads"
description: |-
  Reads Bitbucket downloads via the Bitbucket Cloud API.
---

# bitbucket_downloads (Data Source)

Reads Bitbucket downloads via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/downloads/{filename}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-downloads/#api-repositories-workspace-repo-slug-downloads-filename-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/downloads` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-downloads/#api-repositories-workspace-repo-slug-downloads-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_downloads" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "downloads_response" {
  value = data.bitbucket_downloads.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `filename` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
