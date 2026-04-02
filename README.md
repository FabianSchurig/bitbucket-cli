# bitbucket-cli

A low-maintenance CLI and MCP server for [Bitbucket Cloud](https://bitbucket.org/). Most code is **auto-generated** from the live Bitbucket OpenAPI spec — only a thin hand-written layer ties it together. A daily CI job fetches the latest spec, regenerates the code, and releases a new version if anything changed.

## Why?

Bitbucket Cloud has no official CLI. Managing pull requests through the web UI is slow when you just want to list, approve, merge, or decline from the terminal. This project fills that gap with two binaries:

- **`bb-cli`** — A command-line interface for all Bitbucket Cloud API endpoints.
- **`bb-mcp`** — A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server that exposes all Bitbucket operations as MCP tools.

Both share the same auto-generated foundation:

- **Stays up-to-date automatically** — new API endpoints appear without manual work.
- **Requires near-zero maintenance** — the generic dispatch layer means no per-endpoint glue code.
- **Works everywhere** — Linux, macOS, Windows; install via `go install` or download a binary.
- **Built for AI agents** — designed to be called by coding assistants like GitHub Copilot, Cursor, and similar tools to automate PR workflows: post summaries, add review comments, approve or merge pull requests, and more.

## Architecture

```mermaid
flowchart LR
    A["Bitbucket OpenAPI spec\n(live)"] --> B["enrich_spec.py\n+ operationIds"]
    B --> C["partition_spec.py\nextract paths by group"]
    C --> D["schema/*.yaml"]
    D --> E["oapi-codegen\nmodels.gen.go"]
    D --> F["gen_commands\ncommands/*.gen.go"]
    D --> G["gen_mcptools\nmcptools/*.gen.go"]
    E & F --> H["bb-cli binary"]
    E & G --> I["bb-mcp binary"]
    H --> J["auth · dispatch · output\n(hand-written)"]
    I --> K["MCP handler · dispatch\n(hand-written)"]
    J & K --> L["Shared spec parsing\nscripts/internal/spec"]
```

The architecture uses a shared intermediate representation (`OperationDef`) that
both CLI and MCP generators consume. This design supports future extensions
(e.g., Terraform provider) with minimal additional maintenance.

## Example Usage

List open pull requests:

```bash
bb-cli pr list-pull-requests --workspace myteam --repo-slug myrepo
```

Add a comment on a pull request:

```bash
bb-cli pr create-acomment-on-apull-request \
  --workspace myteam --repo-slug myrepo --pull-request-id 42 \
  --content-raw "Looks good — approved!"
```

List comments with markdown output (useful for AI agents):

```bash
bb-cli pr list-comments-on-apull-request \
  --workspace myteam --repo-slug myrepo --pull-request-id 42 \
  --output markdown
```

Merge a pull request:

```bash
bb-cli pr merge-apull-request \
  --workspace myteam --repo-slug myrepo --pull-request-id 42
```

See all available PR commands:

```bash
bb-cli pr --help
```

## MCP Server

The `bb-mcp` binary is a Model Context Protocol (MCP) server that exposes all Bitbucket API operations as MCP tools. Each command group (pull requests, repositories, pipelines, etc.) is a single tool with an `operation` parameter — this CRUD-combined design matches how Terraform providers work.

### Quick Start

```bash
# Install
go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-mcp@latest

# Set Bitbucket auth
export BITBUCKET_USERNAME=myuser
export BITBUCKET_APP_PASSWORD=ATBBxxxxxxxx

# Run as stdio MCP server (default, for MCP clients like Claude Desktop)
bb-mcp

# Or run as HTTP SSE server
bb-mcp --transport sse --addr :8080
```

### Configuration for Claude Desktop

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "bitbucket": {
      "command": "bb-mcp",
      "env": {
        "BITBUCKET_USERNAME": "myuser",
        "BITBUCKET_APP_PASSWORD": "ATBBxxxxxxxx"
      }
    }
  }
}
```

### Available Tools

Each tool groups related operations with an `operation` parameter:

| Tool | Operations | Description |
|------|-----------|-------------|
| `bitbucket_pr` | 37 | Pull requests: list, create, merge, approve, comments |
| `bitbucket_repos` | 30 | Repositories: list, create, update, permissions |
| `bitbucket_pipelines` | 68 | Pipelines: runs, steps, variables, caches |
| `bitbucket_issues` | 33 | Issues: list, create, update, comments, attachments |
| `bitbucket_workspaces` | 21 | Workspaces: members, permissions, projects |
| ... | | 20 tool groups total, 352+ operations |

## Installation

### Go install (recommended)

```bash
# CLI
go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-cli@latest

# MCP server
go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-mcp@latest
```

Make sure `$(go env GOPATH)/bin` is in your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Shell completion

```bash
# Bash
bb-cli completion bash > /etc/bash_completion.d/bb-cli

# Zsh
bb-cli completion zsh > "${fpath[1]}/_bb-cli"

# Fish
bb-cli completion fish > ~/.config/fish/completions/bb-cli.fish

# PowerShell
bb-cli completion powershell > bb-cli.ps1
```

### Download binary

Download a pre-built binary from the [GitHub Releases](https://github.com/FabianSchurig/bitbucket-cli/releases) page. Archives are available for Linux, macOS, and Windows (amd64/arm64).

### Docker

Two images are published to GHCR on every release:

| Image | Description |
|-------|-------------|
| `ghcr.io/fabianschurig/bitbucket-cli` | CLI (`bb-cli`) |
| `ghcr.io/fabianschurig/bitbucket-mcp` | MCP server (`bb-mcp`) |

Both images use the hardened `dhi.io/golang` base image.

```bash
# CLI
docker pull ghcr.io/fabianschurig/bitbucket-cli:latest

docker run --rm \
  -e BITBUCKET_USERNAME \
  -e BITBUCKET_APP_PASSWORD \
  ghcr.io/fabianschurig/bitbucket-cli pr list-pull-requests \
    --workspace myteam --repo-slug myrepo

# MCP server (stdio, e.g. for Claude Desktop)
docker run --rm -i \
  -e BITBUCKET_USERNAME \
  -e BITBUCKET_APP_PASSWORD \
  ghcr.io/fabianschurig/bitbucket-mcp

# MCP server (SSE over HTTP)
docker run --rm -p 8080:8080 \
  -e BITBUCKET_USERNAME \
  -e BITBUCKET_APP_PASSWORD \
  ghcr.io/fabianschurig/bitbucket-mcp --transport sse --addr :8080
```

#### Building Docker images locally

```bash
# Build CLI image
docker build --target bb-cli -t bb-cli .

# Build MCP server image
docker build --target bb-mcp -t bb-mcp .
```

#### Extending the Dockerfile

Each target is a self-contained stage that installs a binary with `go install`
on the hardened base image. To add a new binary target:

1. Add a new stage that installs the binary with `go install` (use the existing `bb-cli` or `bb-mcp` stages as a template).
2. Add the new target to the build matrix in `.github/workflows/docker.yml` so CI builds and pushes the image automatically.

### Build from source

```bash
git clone https://github.com/FabianSchurig/bitbucket-cli.git
cd bitbucket-cli
go build -o bb-cli ./cmd/...
```

## Contributing

### Dev Container

This repository includes a [Dev Container](https://containers.dev/) configuration that provides a ready-to-use development environment with both **Go** and **Python** pre-installed.

#### Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Visual Studio Code](https://code.visualstudio.com/) with the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers), **or**
- [GitHub Codespaces](https://github.com/features/codespaces)

#### Getting started

1. Clone the repository and open it in VS Code.
2. When prompted, click **Reopen in Container**, or run the command **Dev Containers: Reopen in Container** from the Command Palette (`Ctrl+Shift+P` / `Cmd+Shift+P`).
3. VS Code will build the container and install all required tools. This may take a few minutes on the first run.

After the container starts you will have:

- **Go 1.25** – for building and testing the CLI.
- **Python 3.12** – for running the helper scripts under `scripts/`.
- Go and Python VS Code extensions pre-installed and configured.

#### Running the CLI

```bash
go run ./cmd/bb-cli --help
```

#### Running the scripts

```bash
python3 scripts/enrich_spec.py <input.json> <output.json>
```