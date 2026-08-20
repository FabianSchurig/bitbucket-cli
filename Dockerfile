# ============================================================================
# Dockerfile for bitbucket-cli
#
# Uses a two-stage build per target:
#   1. golang:1 (official, has git) — builds a static binary with `go install`
#   2. scratch — receives only the compiled binary and CA bundle
#
# The runtime image is intentionally package-free to minimize the attack
# surface; the builder stage handles all compilation.
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
RUN CGO_ENABLED=0 GOPROXY=direct go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-cli@${VERSION}

# --- bb-cli: minimal runtime ---
FROM scratch AS bb-cli

COPY --from=build-bb-cli /go/bin/bb-cli /usr/local/bin/bb-cli
COPY --from=build-bb-cli /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER 65532:65532

ENTRYPOINT ["/usr/local/bin/bb-cli"]

# --- bb-mcp: build stage ---
FROM golang:1 AS build-bb-mcp

ARG VERSION=latest

# Use GOPROXY=direct so Go fetches directly from GitHub, bypassing the slow proxy cache
RUN CGO_ENABLED=0 GOPROXY=direct go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-mcp@${VERSION}

# --- bb-mcp: minimal runtime ---
FROM scratch AS bb-mcp

COPY --from=build-bb-mcp /go/bin/bb-mcp /usr/local/bin/bb-mcp
COPY --from=build-bb-mcp /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER 65532:65532

LABEL io.modelcontextprotocol.server.name="io.github.FabianSchurig/bitbucket-mcp"

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/bb-mcp"]
