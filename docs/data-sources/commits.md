---
page_title: "bitbucket_commits Data Source - bitbucket"
subcategory: "Commits"
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
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "commits_response" {
  value = data.bitbucket_commits.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `commit` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `author` (Object) author
  Nested schema:
  - `raw` (String) The raw author value from the repository. This may be the only value available if the author does not match a user in Bitbucket.

- `committer` (Object) committer
  Nested schema:
  - `raw` (String) The raw committer value from the repository. This may be the only value available if the committer does not match a user in Bitbucket.

- `date` (String) date
- `hash` (String) hash
- `message` (String) message
- `parents` (String) parents (JSON array)
- `participants` (List of Object) participants
  Nested schema:
  - `role` (String) [PARTICIPANT, REVIEWER]
  - `approved` (String) approved
  - `state` (String) [approved, changes_requested, <nil>]
  - `participated_on` (String) The ISO8601 timestamp of the participant's action. For approvers, this is the time of their approval. For commenters and pull request reviewers who are not approvers, this is the time they last commented, or null if they have not commented.

- `repository` (Object) repository
  Nested schema:
  - `scm` (String) [git]
  - `created_on` (String) created_on
  - `updated_on` (String) updated_on
  - `language` (String) language
  - `has_wiki` (String) 
  - `uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user.
  - `is_private` (String) is_private
  - `description` (String) description
  - `size` (String) size
  - `has_issues` (String) 
  - `fork_policy` (String) 
  - `full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs.
  - `name` (String) name

- `summary` (Object) summary
  Nested schema:
  - `markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
  - `raw` (String) The text as it was typed by a user.

