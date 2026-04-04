---
page_title: "bitbucket_commits Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket commits via the Bitbucket Cloud API.
---

# bitbucket_commits (Resource)

Manages Bitbucket commits via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_commits" "example" {
  commit = "abc123def"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `commit` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `author_raw` (String) The raw author value from the repository. This may be the only value available if the author does not match a user in...
- `committer_raw` (String) The raw committer value from the repository. This may be the only value available if the committer does not match a u...
- `date` (String) date
- `hash` (String) hash
- `message` (String) message
- `summary_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `summary_raw` (String) The text as it was typed by a user.
