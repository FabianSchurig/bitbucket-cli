# bb-cli usage guide

`bb-cli` is the best entry point for software engineers and computer scientists who want direct Bitbucket Cloud access from the terminal.

## Install

```bash
go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-cli@latest
```

You can also download binaries from the [GitHub Releases](https://github.com/FabianSchurig/bitbucket-cli/releases) page.

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

## Mental model

- Commands are grouped by Bitbucket API area such as `pr`, `repos`, `pipelines`, or `issues`.
- Generated command names stay close to the Bitbucket API operation names.
- `--output table|json|id` controls rendering.
- Pagination follows Bitbucket `next` links automatically unless you pass `--all=false`.
- Nested request body fields become flattened flags such as `source.branch.name` → `--source-branch-name`.

## Common workflows

List pull requests:

```bash
bb-cli pr list-pull-requests --workspace myteam --repo-slug myrepo
```

Show machine-readable output:

```bash
bb-cli repos list-repositories-for-auser --workspace myteam --output json
```

Merge a pull request:

```bash
bb-cli pr merge-apull-request --workspace myteam --repo-slug myrepo --pull-request-id 42
```

## Discover commands quickly

```bash
bb-cli --help
bb-cli pr --help
bb-cli pr list-pull-requests --help
```

If you know the Bitbucket API area but not the exact command name, start with the group help first.

## When to use the CLI

Use `bb-cli` when you want:

- fast terminal access to Bitbucket Cloud
- scripts and shell automation
- easy inspection with `json` output
- direct control without adding Terraform state or an MCP host

## Related links

- [Canonical repository](https://github.com/FabianSchurig/bitbucket-cli)
- [MCP guide](./mcp.md)
- [Terraform provider docs](./index.md)
