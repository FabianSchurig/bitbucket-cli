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
- **golangci-lint** — for linting (`go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`)

## Build & Test

```bash
go build ./...              # build
go test -race ./...         # test with race detector
go run ./cmd/bb-cli --help  # run locally
```

Or use the Makefile:

```bash
make build       # build all binaries
make test        # run tests with race detector
make lint        # run golangci-lint
make fmt         # format all Go files
make fmt-check   # check formatting without modifying files
make vet         # run go vet
make check       # run all checks (build + vet + lint + fmt-check + test)
make generate    # regenerate all code from schemas
make clean       # remove build artifacts
```

## Linting & Static Analysis

This project uses [golangci-lint](https://golangci-lint.run/) with the following linters enabled:

- **errcheck** — ensures all errors are handled
- **govet** — built-in Go static analysis
- **staticcheck** — industry-standard static analyzer (includes gosimple)
- **ineffassign** — detects unused assignments
- **unused** — detects unused code
- **gofmt** — enforces standard Go formatting

Run locally:

```bash
make lint        # or: golangci-lint run ./...
```

Generated files (`*.gen.go`, `internal/generated/`) are excluded from linting.

CI will fail if any linter reports issues.

## Error Handling Guidelines

This project follows strict error handling practices:

1. **Never ignore errors.** All error returns must be checked. Use `_ =` only when the error is genuinely unrecoverable (e.g., `fmt.Fprintf` to stdout in output rendering).
2. **Wrap errors with context** using `fmt.Errorf("...: %w", err)` so callers can understand the failure chain.
3. **Prefer returning errors over panicking.** Only use `panic()` in truly unrecoverable states (e.g., `mustMarshal` for compile-time-known-safe JSON).
4. **Use `recover()` only in well-defined wrapper layers**, not as a general error handling strategy.

The `errcheck` linter enforces these rules in CI.

## Project Structure

```
cmd/bb-cli/main.go          # Entry point, root Cobra command (hand-written)
cmd/bb-mcp/main.go          # Entry point, MCP server (hand-written)
cmd/bb-tf/main.go            # Entry point, Terraform provider (hand-written)
internal/client/             # HTTP client + auth (hand-written)
internal/handlers/           # Generic HTTP dispatcher (hand-written)
internal/output/             # Table/JSON/ID rendering (hand-written)
internal/mcptools/handler.go # MCP tool handler (hand-written)
internal/tfprovider/provider.go   # Terraform provider (hand-written)
internal/tfprovider/resource.go   # Generic CRUD resource (hand-written)
internal/tfprovider/datasource.go # Generic data source (hand-written)
internal/tfprovider/helpers.go    # Shared types + CRUD mapping (hand-written)
internal/commands/*.gen.go   # ⚠️ GENERATED — do not edit
internal/mcptools/*.gen.go   # ⚠️ GENERATED — do not edit
internal/tfprovider/*.gen.go # ⚠️ GENERATED — do not edit
internal/generated/models.gen.go    # ⚠️ GENERATED — do not edit
schema/*-schema.yaml                # ⚠️ GENERATED — do not edit
scripts/internal/spec/       # Shared schema parsing (hand-written)
scripts/gen_commands/         # CLI command generator
scripts/gen_mcptools/         # MCP tool generator
scripts/gen_terraform/        # Terraform resource generator
scripts/                     # Schema enrichment/partition (Python)
```

## Important: Generated Files

The following files are **auto-generated** and must **never be edited by hand**:

- `internal/commands/*.gen.go`
- `internal/mcptools/*.gen.go`
- `internal/tfprovider/*.gen.go`
- `internal/generated/models.gen.go`
- `schema/*-schema.yaml`

Changes will be overwritten by the next CI run. Instead, fix the source:

| Problem | Fix in |
|---------|--------|
| Wrong command flags/descriptions | `scripts/gen_commands/main.go` |
| Wrong MCP tool definitions | `scripts/gen_mcptools/main.go` |
| Wrong Terraform resource definitions | `scripts/gen_terraform/main.go` |
| Wrong model types | `oapi-codegen.yaml` or `scripts/partition_spec.py` |
| Missing/wrong endpoints | `scripts/enrich_spec.py` or `scripts/partition_spec.py` |
| Shared schema logic | `scripts/internal/spec/` |

Generated code is automatically formatted using `go/format` during generation.

## Pull Request Guidelines

1. **Keep changes minimal.** This project values simplicity — avoid adding features or abstractions beyond what's needed.
2. **Don't edit generated files.** Fix the generator or schema script instead.
3. **Include tests** for any new hand-written code in `internal/`.
4. **Run `make check`** (or `go build ./... && go test -race ./... && golangci-lint run ./...`) before submitting.
5. **One concern per PR.** Small, focused PRs are easier to review.
6. **Handle all errors.** CI will reject code with unchecked error returns.
7. **Keep code formatted.** Run `gofmt -w .` or `make fmt` before committing.

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

# Terraform resources
go run scripts/gen_terraform/main.go schema/pr-schema.yaml internal/tfprovider/pr.gen.go

# Or regenerate everything:
make generate
```

## Security

To report a security vulnerability, see [SECURITY.md](SECURITY.md).
