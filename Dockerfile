# ============================================================================
# Dockerfile for bitbucket-cli
#
# Uses a two-stage build per target:
#   1. golang:1 (official, has git) — builds the binary with `go install`
#   2. dhi.io/golang:1 (hardened runtime) — receives only the compiled binary
#
# The hardened base image has no package manager, so we cannot install git
# into it directly; the builder stage handles all compilation.
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
#   To add a new binary target, add a new build+runtime stage pair, then add
#   the target to the build matrix in .github/workflows/docker.yml.
# ============================================================================

# --- bb-cli: build stage ---
FROM golang:1 AS build-bb-cli

ARG VERSION=latest

# Use GOPROXY=direct so Go fetches directly from GitHub, bypassing the slow proxy cache
RUN GOPROXY=direct go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-cli@${VERSION}

# --- bb-cli: hardened runtime ---
FROM dhi.io/golang:1 AS bb-cli

COPY --from=build-bb-cli /go/bin/bb-cli /usr/local/bin/bb-cli

ENTRYPOINT ["bb-cli"]

# --- bb-mcp: build stage ---
FROM golang:1 AS build-bb-mcp

ARG VERSION=latest

# Use GOPROXY=direct so Go fetches directly from GitHub, bypassing the slow proxy cache
RUN GOPROXY=direct go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-mcp@${VERSION}

# --- bb-mcp: hardened runtime ---
FROM dhi.io/golang:1 AS bb-mcp

COPY --from=build-bb-mcp /go/bin/bb-mcp /usr/local/bin/bb-mcp

LABEL io.modelcontextprotocol.server.name="io.github.FabianSchurig/bitbucket-mcp"

EXPOSE 8080
ENTRYPOINT ["bb-mcp"]
