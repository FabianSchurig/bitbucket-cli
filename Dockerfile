# ============================================================================
# Dockerfile for bitbucket-cli
#
# Uses the hardened dhi.io/golang base image and builds binaries from source
# using `go install` (GOPROXY=direct bypasses the module proxy to avoid cache
# lag for newly released versions).
#
# Targets:
#   bb-cli  — Bitbucket CLI
#   bb-mcp  — Bitbucket MCP server (Docker default: last stage)
#
# Build examples:
#   docker build -t bb-mcp .                 # uses default (bb-mcp)
#   docker build --target bb-cli -t bb-cli .
#   docker build --target bb-mcp -t bb-mcp .
#   docker build --build-arg VERSION=v1.0.0 -t bb-mcp .  # pin a version
# Extending this Dockerfile:
#   To add a new binary target, add a new stage that installs the binary
#   with `go install` (use the existing stages as a template), then add
#   the target to the build matrix in .github/workflows/docker.yml.
# ============================================================================

# --- bb-cli: hardened image for the Bitbucket CLI ---
FROM dhi.io/golang:1 AS bb-cli

ARG VERSION=latest

# Use GOPROXY=direct so Go fetches directly from GitHub, bypassing the slow proxy cache
RUN GOPROXY=direct go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-cli@${VERSION}

ENTRYPOINT ["bb-cli"]

# --- bb-mcp: hardened image for the Bitbucket MCP server ---
FROM dhi.io/golang:1 AS bb-mcp

ARG VERSION=latest

# Use GOPROXY=direct so Go fetches directly from GitHub, bypassing the slow proxy cache
RUN GOPROXY=direct go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-mcp@${VERSION}

LABEL io.modelcontextprotocol.server.name="io.github.fabianschurig/bitbucket-mcp"

EXPOSE 8080
ENTRYPOINT ["bb-mcp"]
