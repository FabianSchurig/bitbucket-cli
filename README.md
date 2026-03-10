# bitbucket-cli

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
go run ./cmd/... --help
```

#### Running the scripts

```bash
python3 scripts/enrich_spec.py <input.json> <output.json>
```