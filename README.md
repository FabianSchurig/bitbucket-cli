# bitbucket-cli

A low-maintenance CLI for [Bitbucket Cloud](https://bitbucket.org/) pull requests. Most code is **auto-generated** from the live Bitbucket OpenAPI spec — only a thin hand-written layer ties it together. A daily CI job fetches the latest spec, regenerates the code, and releases a new version if anything changed.

## Why?

Bitbucket Cloud has no official CLI. Managing pull requests through the web UI is slow when you just want to list, approve, merge, or decline from the terminal. This project fills that gap with a single binary that:

- **Stays up-to-date automatically** — new API endpoints appear without manual work.
- **Requires near-zero maintenance** — the generic dispatch layer means no per-endpoint glue code.
- **Works everywhere** — Linux, macOS, Windows; install via `go install` or download a binary.
- **Built for AI agents** — designed to be called by coding assistants like GitHub Copilot, Cursor, and similar tools to automate PR workflows: post summaries, add review comments, approve or merge pull requests, and more.

## Architecture

```mermaid
flowchart LR
    A["Bitbucket OpenAPI spec\n(live)"] --> B["enrich_spec.py\n+ operationIds"]
    B --> C["partition_spec.py\nextract PR paths"]
    C --> D["pr-schema.yaml"]
    D --> E["oapi-codegen\nmodels.gen.go"]
    D --> F["gen_commands\ncommands.gen.go"]
    E & F --> G["bb-cli binary"]
    G --> H["auth · dispatch · output\n(hand-written)"]
```

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

## Installation

### Go install (recommended)

```bash
go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-cli@latest
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

```bash
docker pull ghcr.io/fabianschurig/bitbucket-cli:latest

docker run --rm \
  -e BITBUCKET_USERNAME \
  -e BITBUCKET_APP_PASSWORD \
  ghcr.io/fabianschurig/bitbucket-cli pr list-pull-requests \
    --workspace myteam --repo-slug myrepo
```

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

- **Go 1.24** – for building and testing the CLI.
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