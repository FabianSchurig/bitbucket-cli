---
description: "Use when modifying the HTTP dispatch handler, auth client, or output formatter. Covers hand-written Go runtime code in internal/client, internal/handlers, and internal/output."
applyTo: ["internal/client/**", "internal/handlers/**", "internal/output/**"]
---
# Hand-Written Runtime Code

## Design goals

- **Keep it generic**: `handlers.Dispatch()` handles all API operations uniformly. Avoid per-endpoint branching.
- **Minimize surface area**: The less hand-written code, the less maintenance. Only add code that benefits all operations.
- **Resty for HTTP**: Use `go-resty/resty/v2` — do not introduce additional HTTP libraries.
- **Cobra for CLI**: Use `spf13/cobra` — the root command is in `cmd/main.go`, subcommands are generated.

## Module responsibilities

| Module | Responsibility |
|--------|---------------|
| `internal/client/auth.go` | Resty client factory with Basic/Bearer auth from env vars |
| `internal/handlers/dispatch.go` | Generic HTTP dispatcher: URL templating, query params, pagination, body building |
| `internal/output/format.go` | Render responses as `table`, `json`, or `id` format |

## Pagination

Bitbucket uses cursor-based pagination (`values` array + `next` URL). The `--all` flag triggers automatic page-following in `Dispatch()`.

## Auth priority

1. `BITBUCKET_USERNAME` + `BITBUCKET_APP_PASSWORD` → HTTP Basic
2. `BITBUCKET_TOKEN` alone → Bearer (OAuth2)

## Testing

```bash
go test ./internal/...
```
