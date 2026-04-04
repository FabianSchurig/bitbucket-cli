# Bitbucket CLI — Project Guidelines

## Overview

A low-maintenance CLI and MCP server for Bitbucket Cloud. Most code is **auto-generated** from the live Bitbucket OpenAPI spec — only a thin hand-written layer ties it together.

## Architecture

```
Bitbucket OpenAPI spec (live)
  → scripts/enrich_spec.py        # inject operationIds
  → scripts/partition_spec.py     # extract paths by group, resolve $refs
  → schema/*-schema.yaml          # self-contained OpenAPI specs (one per group)
  → oapi-codegen                  # internal/generated/models.gen.go
  → scripts/gen_commands/main.go  # internal/commands/*.gen.go (CLI)
  → scripts/gen_mcptools/main.go  # internal/mcptools/*.gen.go (MCP)
```

Shared code generation logic lives in `scripts/internal/spec/` (schema types,
body field resolution, operation building).

Hand-written code lives in:
- `cmd/bb-cli/main.go` — CLI entry point, root Cobra command
- `cmd/bb-mcp/main.go` — MCP server entry point (stdio + SSE transports)
- `internal/client/auth.go` — Resty client + auth (Basic or Bearer)
- `internal/handlers/dispatch.go` — generic HTTP dispatcher with pagination
- `internal/output/format.go` — table / json / id rendering
- `internal/mcptools/handler.go` — generic MCP tool handler (CRUD-combined dispatch)

## Critical Rules

1. **Never hand-edit generated files.** `internal/commands/*.gen.go`, `internal/mcptools/*.gen.go`, and `internal/generated/models.gen.go` are produced by the pipeline. Fix the generator or schema instead.
2. **Minimize hand-written code.** The design goal is near-zero maintenance — new Bitbucket endpoints arrive automatically via schema sync.
3. **Keep the dispatch generic.** `handlers.Dispatch()` and `handlers.DispatchRaw()` handle all operations uniformly. Avoid per-endpoint special cases unless absolutely unavoidable.
4. **Shared schema parsing.** `scripts/internal/spec/` contains shared types and helpers used by both `gen_commands` and `gen_mcptools`. Changes here affect both generators.

## Environment

- **Dev container**: Ubuntu base, Go 1.24, Python 3.12
- **Go deps**: `cobra` (CLI), `resty/v2` (HTTP), `yaml.v3` (schema parsing), `go-sdk/mcp` (MCP server)
- **Python deps**: `pyyaml` (schema scripts only)
- **Code gen**: `oapi-codegen` for models

## Build & Test

```bash
go build ./...          # build
go test ./...           # test
go run ./cmd/bb-cli --help # run CLI locally
go run ./cmd/bb-mcp       # run MCP server locally
```

## Code Generation (manual)

```bash
python3 scripts/enrich_spec.py <raw-spec.json> <enriched.json>
python3 scripts/partition_spec.py <enriched.json> schema/ --all
oapi-codegen --config oapi-codegen.yaml schema/pr-schema.yaml
go run scripts/gen_commands/main.go schema/pr-schema.yaml internal/commands/commands.gen.go
go run scripts/gen_mcptools/main.go schema/pr-schema.yaml internal/mcptools/pr.gen.go
```

## Conventions

- **Auth**: `BITBUCKET_USERNAME` + `BITBUCKET_API_TOKEN` (Basic) or `BITBUCKET_TOKEN` (Bearer)
- **Flags**: path params → `--workspace`, `--repo-slug`; body fields → flattened with dots (`--source-branch` maps to `source.branch`)
- **Output**: `--output table|json|id`; table is default
- **Pagination**: `--all` flag (default: **true**) auto-follows cursor-based `next` links; pass `--all=false` for a single page
- **Versioning**: semver `v1.x.y` — compatible with `go install` / pkg.go.dev; patch is auto-bumped by CI on schema changes
- **CI**: Daily cron at 03:00 UTC fetches spec, regenerates, tests, and releases if changed
