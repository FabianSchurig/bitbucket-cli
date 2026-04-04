---
page_title: "bitbucket_commits Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket commits via the Bitbucket Cloud API.
---

# bitbucket_commits (Data Source)

Reads Bitbucket commits via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/commit/{commit}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commit-commit-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/commits` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commits-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

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
- `author_raw` (String) The raw author value from the repository. This may be the only value available if the author does not match a user in Bitbucket.
- `committer_raw` (String) The raw committer value from the repository. This may be the only value available if the committer does not match a user in Bitbucket.
- `date` (String) date
- `hash` (String) hash
- `message` (String) message
- `summary_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `summary_raw` (String) The text as it was typed by a user.
