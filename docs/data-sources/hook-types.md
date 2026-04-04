---
page_title: "bitbucket_hook_types Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket hook-types via the Bitbucket Cloud API.
---

# bitbucket_hook_types (Data Source)

Reads Bitbucket hook-types via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/hook_events` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-hook-events-get) |
| List | `GET` | `/hook_events/{subject_type}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-hook-events-subject-type-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | — |
| List | — |

## Example Usage

```hcl
data "bitbucket_hook_types" "example" {
}

output "hook_types_response" {
  value = data.bitbucket_hook_types.example.api_response
}
```

## Schema

### Required

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `repository_events_href` (String) repository.events.href
- `repository_events_name` (String) repository.events.name
- `workspace_events_href` (String) workspace.events.href
- `workspace_events_name` (String) workspace.events.name
