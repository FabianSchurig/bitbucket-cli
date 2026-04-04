---
page_title: "bitbucket_commit_file Data Source - bitbucket"
subcategory: "Repositories"
description: |-
  Reads Bitbucket commit-file via the Bitbucket Cloud API.
---

# bitbucket_commit_file (Data Source)

Reads Bitbucket commit-file via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/src/{commit}/{path}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-src-commit-path-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_commit_file" "example" {
  commit = "abc123def"
  path = "README.md"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "commit_file_response" {
  value = data.bitbucket_commit_file.example.api_response
}
```

## Schema

### Required
- `commit` (String) Path parameter.
- `path` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `commit_date` (String) commit.date
- `commit_hash` (String) commit.hash
- `commit_message` (String) commit.message
- `commit_parents` (String) parents (JSON array)
- `commit_participants` (List of Object) participants
  Nested schema:
  - `role` (String) [PARTICIPANT, REVIEWER]
  - `approved` (String) approved
  - `state` (String) [approved, changes_requested, <nil>]
  - `participated_on` (String) The ISO8601 timestamp of the participant's action. For approvers, this is the time of their approval. For commenters and pull request reviewers who are not approvers, this is the time they last commented, or null if they have not commented.

- `commit_summary_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `commit_summary_raw` (String) The text as it was typed by a user.
- `type` (String) type
