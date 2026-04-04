# bb-mcp usage guide

`bb-mcp` is the best entry point for AI agents and MCP-compatible clients that need Bitbucket Cloud tools.

## Install

```bash
go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-mcp@latest
```

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

## Example Claude Desktop configuration

```json
{
  "mcpServers": {
    "bitbucket": {
      "command": "bb-mcp",
      "env": {
        "BITBUCKET_USERNAME": "your-email@example.com",
        "BITBUCKET_TOKEN": "your-api-token"
      }
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

## Related links

- [Canonical repository](https://github.com/FabianSchurig/bitbucket-cli)
- [CLI guide](./cli.md)
- [Terraform provider docs](./index.md)
