---
description: "Use when modifying the command generator, schema enrichment, or partitioning scripts. Covers the code generation pipeline that produces CLI commands and models from the Bitbucket OpenAPI spec."
applyTo: ["scripts/**", "oapi-codegen.yaml"]
---
# Code Generation Pipeline

## Pipeline stages

1. **`scripts/enrich_spec.py`** — Injects `operationId` into raw Bitbucket OpenAPI spec (summary→camelCase with deduplication)
2. **`scripts/partition_spec.py`** — Extracts PR-related paths, recursively resolves `$ref`s into a self-contained schema
3. **`oapi-codegen`** — Generates Go model types from `schema/pr-schema.yaml`
4. **`scripts/gen_commands/main.go`** — Generates Cobra commands wired to `handlers.Dispatch()`

## Key principles

- **Generic dispatch**: Generated commands call `handlers.Dispatch()` uniformly — no per-endpoint special cases
- **Flat flags**: Nested request body fields are flattened into CLI flags (e.g., `source.branch` → `--source-branch`)
- **Minimal hand-written code**: The generator should handle all boilerplate; only truly generic behavior belongs in hand-written code
- **Schema scripts are Python 3.12**: Only dependency is `pyyaml`; keep scripts simple and dependency-light
- **Go generator uses `yaml.v3`**: Parses the OpenAPI schema directly without third-party OpenAPI libraries

## Testing changes

After modifying any generator or schema script:
```bash
# Regenerate and verify
go run scripts/gen_commands/main.go schema/pr-schema.yaml internal/commands/commands.gen.go
go build ./...
go test ./...
```
