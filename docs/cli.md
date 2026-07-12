# bb-cli usage guide

`bb-cli` is the best entry point for software engineers and computer scientists who want direct Bitbucket Cloud access from the terminal.

## Install

### Homebrew (macOS and Linux)

To install via our custom tap, tap the repository and trust the formula to reduce supply chain risks:

```bash
brew tap FabianSchurig/tap
brew trust --formula FabianSchurig/tap/bitbucket-cli
brew install bitbucket-cli
```

> [!NOTE]
> **macOS "Unidentified Developer" Warning**
> Because our pre-built binaries are not signed/notarized with a paid Apple Developer ID, macOS Gatekeeper will block execution with a warning when installed via Homebrew.
> 
> You can bypass this warning by running:
> ```bash
> xattr -d com.apple.quarantine $(which bb-cli)
> ```
> *Alternatively, you can choose **Allow Anyway** under System Settings > Privacy & Security.*
> *In the future, we may obtain a paid Apple Developer Membership to sign and notarize the pre-compiled binaries natively.*

### Install script

Download and install the latest release automatically:

```bash
curl -fsSL https://raw.githubusercontent.com/FabianSchurig/bitbucket-cli/main/install.sh | sh
```

Install a specific version or binary:

```bash
# Install a specific version
curl -fsSL https://raw.githubusercontent.com/FabianSchurig/bitbucket-cli/main/install.sh | sh -s -- --version v1.2.3

# Install bb-mcp instead of bb-cli
curl -fsSL https://raw.githubusercontent.com/FabianSchurig/bitbucket-cli/main/install.sh | sh -s -- --binary bb-mcp

# Install both bb-cli and bb-mcp
curl -fsSL https://raw.githubusercontent.com/FabianSchurig/bitbucket-cli/main/install.sh | sh -s -- --binary all

# Install to a custom directory
curl -fsSL https://raw.githubusercontent.com/FabianSchurig/bitbucket-cli/main/install.sh | sh -s -- --install-dir ~/.local/bin
```

### Go install

If you are on macOS and want to bypass the Gatekeeper warning completely without running `xattr`, you can compile the binary locally using `go install`:

```bash
go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-cli@latest
```

### Docker

Multi-arch container images are published to GHCR on every release:

```bash
docker pull ghcr.io/fabianschurig/bitbucket-cli:latest
docker run -e BITBUCKET_USERNAME -e BITBUCKET_TOKEN ghcr.io/fabianschurig/bitbucket-cli:latest --help
```

### Download binaries

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

## Workspace and repository inference

When you run a command that requires `--workspace` or `--repo-slug`, `bb-cli` can infer them automatically so you don't have to type them every time.

**Precedence** (highest to lowest):

1. Explicit CLI flags (`--workspace`, `--repo-slug`)
2. Environment variables (`BITBUCKET_WORKSPACE`, `BITBUCKET_REPO_SLUG`)
3. Git remote URL of the current directory's `origin` remote

Supported remote URL formats:

- SSH: `git@bitbucket.org:workspace/repo.git`
- HTTPS: `https://bitbucket.org/workspace/repo.git`
- HTTPS with user: `https://user@bitbucket.org/workspace/repo.git`

This means that inside a cloned Bitbucket repository you can simply run:

```bash
bb-cli pr list-pull-requests
```

instead of:

```bash
bb-cli pr list-pull-requests --workspace myteam --repo-slug myrepo
```

If inference fails (not inside a git repo, or the remote is not on Bitbucket), the command reports the usual `--workspace is required` / `--repo-slug is required` error.

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
