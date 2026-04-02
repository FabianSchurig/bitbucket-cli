# Contributing to bitbucket-cli

Thanks for your interest in contributing! This project is designed for **low maintenance** — most code is auto-generated from the Bitbucket OpenAPI spec. Contributions that align with this philosophy are welcome.

## Getting Started

### Dev Container (recommended)

This repository includes a [Dev Container](https://containers.dev/) configuration:

1. Clone the repo and open it in VS Code.
2. Click **Reopen in Container** (or use GitHub Codespaces).
3. The container provides **Go 1.24** and **Python 3.12** with all tools pre-configured.

### Manual Setup

- **Go 1.24+** — for building and testing the CLI
- **Python 3.12+** with `pyyaml` — for running the schema scripts

## Build & Test

```bash
go build ./...          # build
go test ./...           # test
go run ./cmd/bb-cli --help  # run locally
```

## Project Structure

```
cmd/bb-cli/main.go          # Entry point, root Cobra command (hand-written)
cmd/bb-mcp/main.go          # Entry point, MCP server (hand-written)
internal/client/             # HTTP client + auth (hand-written)
internal/handlers/           # Generic HTTP dispatcher (hand-written)
internal/output/             # Table/JSON/ID rendering (hand-written)
internal/mcptools/handler.go # MCP tool handler (hand-written)
internal/commands/*.gen.go   # ⚠️ GENERATED — do not edit
internal/mcptools/*.gen.go   # ⚠️ GENERATED — do not edit
internal/generated/models.gen.go    # ⚠️ GENERATED — do not edit
schema/*-schema.yaml                # ⚠️ GENERATED — do not edit
scripts/internal/spec/       # Shared schema parsing (hand-written)
scripts/gen_commands/         # CLI command generator
scripts/gen_mcptools/         # MCP tool generator
scripts/                     # Schema enrichment/partition (Python)
```

## Important: Generated Files

The following files are **auto-generated** and must **never be edited by hand**:

- `internal/commands/*.gen.go`
- `internal/mcptools/*.gen.go`
- `internal/generated/models.gen.go`
- `schema/*-schema.yaml`

Changes will be overwritten by the next CI run. Instead, fix the source:

| Problem | Fix in |
|---------|--------|
| Wrong command flags/descriptions | `scripts/gen_commands/main.go` |
| Wrong MCP tool definitions | `scripts/gen_mcptools/main.go` |
| Wrong model types | `oapi-codegen.yaml` or `scripts/partition_spec.py` |
| Missing/wrong endpoints | `scripts/enrich_spec.py` or `scripts/partition_spec.py` |
| Shared schema logic | `scripts/internal/spec/` |

## Pull Request Guidelines

1. **Keep changes minimal.** This project values simplicity — avoid adding features or abstractions beyond what's needed.
2. **Don't edit generated files.** Fix the generator or schema script instead.
3. **Include tests** for any new hand-written code in `internal/`.
4. **Run `go build ./...` and `go test ./...`** before submitting.
5. **One concern per PR.** Small, focused PRs are easier to review.

## Code Generation (for reference)

The full pipeline, typically run by CI:

```bash
python3 scripts/enrich_spec.py <raw-spec.json> <enriched.json>
python3 scripts/partition_spec.py <enriched.json> schema/ --all
oapi-codegen --config oapi-codegen.yaml schema/pr-schema.yaml

# CLI commands
go run scripts/gen_commands/main.go schema/pr-schema.yaml internal/commands/commands.gen.go

# MCP tools
go run scripts/gen_mcptools/main.go schema/pr-schema.yaml internal/mcptools/pr.gen.go
```

## Security

To report a security vulnerability, see [SECURITY.md](SECURITY.md).
