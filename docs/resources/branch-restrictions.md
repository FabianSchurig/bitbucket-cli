---
page_title: "bitbucket_branch_restrictions Resource - bitbucket"
subcategory: "Branch Restrictions"
description: |-
  Manages Bitbucket branch-restrictions via the Bitbucket Cloud API.
---

# bitbucket_branch_restrictions (Resource)

Manages Bitbucket branch-restrictions via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/branch-restrictions` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branch-restrictions/#api-repositories-workspace-repo-slug-branch-restrictions-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/branch-restrictions/{id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branch-restrictions/#api-repositories-workspace-repo-slug-branch-restrictions-id-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/branch-restrictions/{id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branch-restrictions/#api-repositories-workspace-repo-slug-branch-restrictions-id-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/branch-restrictions/{id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branch-restrictions/#api-repositories-workspace-repo-slug-branch-restrictions-id-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/branch-restrictions` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branch-restrictions/#api-repositories-workspace-repo-slug-branch-restrictions-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `admin:repository:bitbucket` |
| Read | `admin:repository:bitbucket` |
| Update | `admin:repository:bitbucket` |
| Delete | `admin:repository:bitbucket` |
| List | `admin:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_branch_restrictions" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
  branch_match_kind = "glob"
  kind = "require_approvals_to_merge"
  pattern = "main"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.
- `branch_match_kind` (String) Indicates how the restriction is matched against a branch. The default is `glob`. [branching_model, glob]
- `kind` (String) The type of restriction that is being applied. [push, delete, force, restrict_merges, require_tasks_to_be_completed, require_approvals_to_merge, require_review_group_approvals_to_merge, require_default_reviewer_approvals_to_merge, require_no_changes_requested, require_passing_builds_to_merge, require_commits_behind, reset_pullrequest_approvals_on_change, smart_reset_pullrequest_approvals, reset_pullrequest_changes_requested_on_change, require_all_dependencies_merged, enforce_merge_checks, allow_auto_merge_when_builds_pass, require_all_comments_resolved]
- `pattern` (String) Apply the restriction to branches that match this pattern. Active when `branch_match_kind` is `glob`. Will be empty when `branch_match_kind` is `branching_model`.

### Optional
- `param_id` (String) Path parameter (auto-populated from API response).
- `branch_type` (String) Apply the restriction to branches of this type. Active when `branch_match_kind` is `branching_model`. The branch type will be calculated using the branching model configured for the repository. [feature, bugfix, release, hotfix, development, production] (also computed from API response)
- `groups` (List of Object) groups (also computed from API response)
  Nested schema:
  - `full_slug` (String) The concatenation of the workspace's slug and the group's slug,
  - `name` (String) name
  - `slug` (String) The "sluggified" version of the group's name. This contains only ASCII

- `users` (List of Object) users (also computed from API response)
  Nested schema:
  - `created_on` (String) created_on
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `value` (String) Value with kind-specific semantics: (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only
- `api_response` (String) The raw JSON response from the Bitbucket API.
