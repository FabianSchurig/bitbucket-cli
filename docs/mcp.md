# bb-mcp usage guide

`bb-mcp` is the best entry point for AI agents and MCP-compatible clients that need Bitbucket Cloud tools.

## Install

### Go install

```bash
go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-mcp@latest
```

### Docker

Multi-arch container images are published to GHCR on every release:

```bash
docker pull ghcr.io/fabianschurig/bitbucket-mcp:latest
docker run -e BITBUCKET_USERNAME -e BITBUCKET_TOKEN -p 8080:8080 ghcr.io/fabianschurig/bitbucket-mcp:latest --transport sse --addr :8080
```

### MCP Registry

`bb-mcp` is published to the [MCP Registry](https://registry.modelcontextprotocol.io) automatically on every release. MCP-compatible clients can discover and install it from the registry.

## Authenticate

API token:

```bash
export BITBUCKET_USERNAME="your-email@example.com"
export BITBUCKET_TOKEN="your-api-token"
```

Workspace or repository access token:

```bash
export BITBUCKET_TOKEN="your-access-token"
```

## Run the server

Default stdio transport:

```bash
bb-mcp
```

HTTP SSE transport:

```bash
bb-mcp --transport sse --addr :8080
```

## How the tools are structured

- Tools are grouped by Bitbucket area such as pull requests, repositories, pipelines, or issues.
- Each tool accepts an `operation` parameter instead of creating one MCP tool per endpoint.
- Parameters map closely to the Bitbucket API, so required path/query/body inputs stay easy to trace.
- The grouped design keeps the MCP surface smaller while still exposing broad API coverage.

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

To override, create an `mcp_config.yaml` next to the binary (see `mcp_config.yaml` in the repo root for a commented template).

## Example VS Code configuration

```json
{
	"servers": {
		"bitbucket-mcp-server": {
			"type": "stdio",
			"command": "bb-mcp",
			"args": ["--config", "${workspaceFolder}/mcp_config.yaml"],
			"envFile": "${workspaceFolder}/.env"
		}
	},
	"inputs": []
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

## Related links

- [Canonical repository](https://github.com/FabianSchurig/bitbucket-cli)
- [CLI guide](./cli.md)
- [Terraform provider docs](./index.md)
