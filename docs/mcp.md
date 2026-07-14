# bb-mcp usage guide

`bb-mcp` is the best entry point for AI agents and MCP-compatible clients that need Bitbucket Cloud tools.

## Install

Choose whichever install path fits your setup:

- **Direct binary (`bb-mcp`)**: Homebrew, install script, or `go install`
- **Containerized**: Docker image from GHCR

### Homebrew (macOS and Linux)

To install via our custom tap, tap the repository and trust the formula to reduce supply chain risks:

```bash
brew tap FabianSchurig/tap
brew trust --formula FabianSchurig/tap/bitbucket-mcp
brew install bitbucket-mcp
```

> [!NOTE]
> **macOS "Unidentified Developer" Warning**
> Because our pre-built binaries are not signed/notarized with a paid Apple Developer ID, macOS Gatekeeper will block execution with a warning when installed via Homebrew.
> 
> You can bypass this warning by running:
> ```bash
> xattr -d com.apple.quarantine $(which bb-mcp)
> ```
> *Alternatively, you can choose **Allow Anyway** under System Settings > Privacy & Security.*
> *In the future, we may obtain a paid Apple Developer Membership to sign and notarize the pre-compiled binaries natively.*

### Install script

```bash
curl -fsSL https://raw.githubusercontent.com/FabianSchurig/bitbucket-cli/main/install.sh | sh -s -- --binary bb-mcp
```

### Go install

If you are on macOS and want to bypass the Gatekeeper warning completely without running `xattr`, you can compile the binary locally using `go install`:

```bash
go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-mcp@latest
```

### Docker

Multi-arch container images are published to GHCR on every release.

Use a pinned tag for reproducible behavior:

```bash
docker pull ghcr.io/fabianschurig/bitbucket-mcp:0.18.5
docker run --rm -i --env-file "${HOME}/.config/bitbucket-mcp.env" -p 8080:8080 ghcr.io/fabianschurig/bitbucket-mcp:0.18.5 --transport sse --addr :8080
```

> [!TIP]
> If credentials are already exported in your shell, `-e BITBUCKET_USERNAME -e BITBUCKET_TOKEN` forwards host environment values to the container.

### MCP Registry

`bb-mcp` is published to the [MCP Registry](https://registry.modelcontextprotocol.io) automatically on every release. MCP-compatible clients can discover and install it from the registry.

## Authenticate

`bb-mcp` supports Atlassian scoped API tokens (recommended), and Bitbucket workspace/repository access tokens.

> Atlassian scoped API tokens use `:bitbucket`-suffixed scopes and are distinct from legacy Bitbucket App Passwords.

### Create a scoped API token

1. Go to <https://id.atlassian.com/manage-profile/security/api-tokens>.
2. Click **Create API token with scopes**.
3. Select **Bitbucket**, then choose scopes you need, for example:
   - `read:repository:bitbucket`, `read:pullrequest:bitbucket`
   - `read:account`, `read:me`
   - add `write:*` scopes only when you need write access
4. Copy the token (usually starts with `ATATT...`). You cannot view it again later.

### Provide credentials (shared env file location)

Use a single env file location for both direct `bb-mcp` and Docker flows:

`$HOME/.config/bitbucket-mcp.env`

Using an API token (username is your **Atlassian account email**):

```bash
export BITBUCKET_USERNAME="your-email@example.com"
export BITBUCKET_TOKEN="your-api-token"
```

Using a workspace or repository access token (no username required):

```bash
export BITBUCKET_TOKEN="your-access-token"
```

Recommended env file content (**unquoted** values):

```env
BITBUCKET_USERNAME=your-email@example.com
BITBUCKET_TOKEN=ATATT...
```

> [!WARNING]
> **Quotes in `.env` files.** The `export` examples above use quotes because your shell strips them. Docker `--env-file` passes values verbatim, so quoted values include the quote characters and can cause `401 Unauthorized`.

## Run the server

If installed directly:

Default stdio transport:

```bash
bb-mcp
```

HTTP SSE transport:

```bash
bb-mcp --transport sse --addr :8080
```

If running with Docker:

```bash
docker run --rm -i --env-file "${HOME}/.config/bitbucket-mcp.env" -p 8080:8080 ghcr.io/fabianschurig/bitbucket-mcp:0.18.5 --transport sse --addr :8080
```

## How the tools are structured

- Tools are grouped by Bitbucket area such as pull requests, repositories, pipelines, or issues.
- Each tool accepts an `operation` parameter instead of creating one MCP tool per endpoint.
- Parameters map closely to the Bitbucket API, so required path/query/body inputs stay easy to trace.
- The grouped design keeps the MCP surface smaller while still exposing broad API coverage.

## Workspace and repository inference

When a tool requires `workspace` or `repo_slug` parameters, `bb-mcp` can infer them automatically from the server's working directory.

**Precedence** (highest to lowest):

1. Explicit tool parameters (`workspace`, `repo_slug`)
2. Environment variables (`BITBUCKET_WORKSPACE`, `BITBUCKET_REPO_SLUG`)
3. Git remote URL of the current directory's `origin` remote

Supported remote URL formats: SSH (`git@bitbucket.org:ws/repo.git`) and HTTPS (`https://bitbucket.org/ws/repo.git`).

This allows agents to call tools without specifying `workspace` and `repo_slug` when the MCP server runs inside a cloned Bitbucket repository.

## Available tools

A complete auto-generated reference — every tool group, every operation, and all active description overrides — lives in [tools-reference.md](./tools-reference.md).

The reference is regenerated automatically from the Bitbucket OpenAPI schemas whenever the schema changes (via `make generate-docs`).

### Quick summary by area

| Tool | Purpose |
|------|---------|
| `bitbucket_pr` | Pull request CRUD, review, comments, tasks, merge |
| `bitbucket_pipelines` | Run, inspect, and debug CI/CD pipelines |
| `bitbucket_repos` | Browse repos, read source files, manage settings |
| `bitbucket_commits` | Commit history, diffs, branch comparisons |
| `bitbucket_refs` | Branches and tags |
| `bitbucket_search` | Full-text code search across all repos |
| `bitbucket_issues` | Issue tracker |
| `bitbucket_commit-statuses` | CI status checks per commit/PR |
| `bitbucket_deployments` | Deployment environment tracking |
| `bitbucket_branch-restrictions` | Branch protection rules |
| `bitbucket_branching-model` | Gitflow-style branching model settings |
| `bitbucket_workspaces` | Workspace membership and settings |
| `bitbucket_projects` | Project organisation within a workspace |
| `bitbucket_hooks` | Repository and workspace webhooks |
| `bitbucket_reports` | Code-quality reports attached to commits |
| `bitbucket_snippets` | Shared code snippets |
| `bitbucket_users` | User profiles and SSH keys |
| `bitbucket_downloads` | Repository file downloads |

## Default configuration

When no `mcp_config.yaml` is present in the working directory the server uses a built-in default that:

- **Allows** `GET`, `POST`, `PUT`, `PATCH` — no `DELETE` operations exposed.
- **Hides** `bitbucket_addon` and `bitbucket_properties` (platform/admin tools).
- **Applies** LLM-optimised descriptions for the eight most important daily tools.

To override, create an `mcp_config.yaml` in your working directory — use [`internal/config/default_mcp_config.yaml`](../internal/config/default_mcp_config.yaml) as a commented starting point.

## Example VS Code configuration (direct `bb-mcp`)

This uses your local `bb-mcp` binary and the shared env file location:

```json
{
	"servers": {
		"bitbucket-mcp-server": {
			"type": "stdio",
			"command": "bb-mcp",
			"args": ["--config", "${workspaceFolder}/mcp_config.yaml"],
			"envFile": "${userHome}/.config/bitbucket-mcp.env"
		}
	},
	"inputs": []
}
```

## Example Docker configuration (`mcp.json`)

Runs the server in a container, reading credentials from the same shared env file location:

```json
{
  "mcpServers": {
    "bitbucket-mcp-server": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "--env-file",
        "${userHome}/.config/bitbucket-mcp.env",
        "ghcr.io/fabianschurig/bitbucket-mcp:0.18.5"
      ]
    }
  }
}
```

## Good use cases for MCP

Use `bb-mcp` when you want an agent to:

- inspect pull requests, comments, pipelines, repositories, or workspace data
- automate review flows and repository operations from an MCP client
- share one Bitbucket integration across multiple agent prompts or tools

## Choosing the right transport

- **stdio**: best for local MCP clients such as Claude Desktop
- **SSE**: best when your MCP client expects an HTTP endpoint

## Troubleshooting

### `401 Unauthorized`

Credentials were rejected. Common causes:

- **Quotes in `--env-file`.** Remove surrounding quotes in `~/.config/bitbucket-mcp.env` values.
- **Wrong username.** `BITBUCKET_USERNAME` must be your Atlassian **email** (not Bitbucket handle).
- **Stale process.** Restart your MCP server/client after editing the env file.
- **Bad token value.** Token expired/revoked, or includes trailing whitespace/newline.

### `401` vs `403`

- `401` = authentication failed (credentials rejected)
- `403` = authenticated, but missing required scope

## Related links

- [Canonical repository](https://github.com/FabianSchurig/bitbucket-cli)
- [CLI guide](./cli.md)
- [Terraform provider docs](./index.md)
