---
page_title: "bitbucket_commits Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket commits via the Bitbucket Cloud API.
---

# bitbucket_commits (Data Source)

Reads Bitbucket commits via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_commits" "example" {
  commit = "abc123def"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "commits_response" {
  value = data.bitbucket_commits.example.api_response
}
```

## Schema

### Required
- `commit` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `author_raw` (String) The raw author value from the repository. This may be the only value available if the author does not match a user in...
- `committer_raw` (String) The raw committer value from the repository. This may be the only value available if the committer does not match a u...
- `date` (String) date
- `hash` (String) hash
- `message` (String) message
- `summary_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `summary_raw` (String) The text as it was typed by a user.
