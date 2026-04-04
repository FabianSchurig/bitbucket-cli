---
page_title: "bitbucket_branch_restrictions Data Source - bitbucket"
subcategory: "Branch Restrictions"
description: |-
  Reads Bitbucket branch-restrictions via the Bitbucket Cloud API.
---

# bitbucket_branch_restrictions (Data Source)

Reads Bitbucket branch-restrictions via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/branch-restrictions/{id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branch-restrictions/#api-repositories-workspace-repo-slug-branch-restrictions-id-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/branch-restrictions` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branch-restrictions/#api-repositories-workspace-repo-slug-branch-restrictions-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `admin:repository:bitbucket` |
| List | `admin:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_branch_restrictions" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "branch_restrictions_response" {
  value = data.bitbucket_branch_restrictions.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `param_id` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `branch_match_kind` (String) Indicates how the restriction is matched against a branch. The default is `glob`. [branching_model, glob]
- `branch_type` (String) Apply the restriction to branches of this type. Active when `branch_match_kind` is `branching_model`. The branch type will be calculated using the branching model configured for the repository. [feature, bugfix, release, hotfix, development, production]
- `groups` (List of Object) groups
  Nested schema:
  - `name` (String) name
  - `slug` (String) The "sluggified" version of the group's name. This contains only ASCII
  - `full_slug` (String) The concatenation of the workspace's slug and the group's slug,

- `kind` (String) The type of restriction that is being applied. [push, delete, force, restrict_merges, require_tasks_to_be_completed, require_approvals_to_merge, require_review_group_approvals_to_merge, require_default_reviewer_approvals_to_merge, require_no_changes_requested, require_passing_builds_to_merge, require_commits_behind, reset_pullrequest_approvals_on_change, smart_reset_pullrequest_approvals, reset_pullrequest_changes_requested_on_change, require_all_dependencies_merged, enforce_merge_checks, allow_auto_merge_when_builds_pass, require_all_comments_resolved]
- `pattern` (String) Apply the restriction to branches that match this pattern. Active when `branch_match_kind` is `glob`. Will be empty when `branch_match_kind` is `branching_model`.
- `users` (List of Object) users
  Nested schema:
  - `created_on` (String) created_on
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `value` (String) Value with kind-specific semantics:
