---
page_title: "bitbucket_tags Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket tags via the Bitbucket Cloud API.
---

# bitbucket_tags (Data Source)

Reads Bitbucket tags via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/refs/tags/{name}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-tags-name-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/refs/tags` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-tags-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_tags" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
  name = "main"
}

output "tags_response" {
  value = data.bitbucket_tags.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.
- `name` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `date` (String) The date that the tag was created, if available
- `message` (String) The message associated with the tag, if available.
- `name` (String) The name of the ref.
- `tagger_raw` (String) The raw author value from the repository. This may be the only value available if the author does not match a user in Bitbucket.
- `type` (String) type
