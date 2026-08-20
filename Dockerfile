# ============================================================================
# Dockerfile for bitbucket-cli
#
# Uses a shared builder plus a runtime stage per target:
#   1. golang:1 — builds static binaries from the checked-out source tree
#   2. scratch — receives only the compiled binary and CA bundle
#
# The runtime image is intentionally package-free to minimize the attack
# surface; the builder stage handles dependency resolution and compilation.
#
# Targets:
#   bb-cli  — Bitbucket CLI
#   bb-mcp  — Bitbucket MCP server (Docker default: last stage)
#
# Build examples:
#   docker build -t bb-mcp .                 # uses default target (bb-mcp)
#   docker build --target bb-cli -t bb-cli .
#   docker build --target bb-mcp -t bb-mcp .
# Extending this Dockerfile:
#   To add a new binary target, add a new build+runtime stage pair, then add
#   the target to the build matrix in .github/workflows/docker.yml.
# ============================================================================

# --- shared build stage ---
FROM golang:1 AS build-base

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

# --- bb-cli: build stage ---
FROM build-base AS build-bb-cli

RUN CGO_ENABLED=0 go build -o /out/bb-cli ./cmd/bb-cli

# --- bb-cli: minimal runtime ---
FROM scratch AS bb-cli

COPY --from=build-bb-cli /out/bb-cli /bb-cli
COPY --from=build-bb-cli /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER 65532:65532

ENTRYPOINT ["/bb-cli"]

# --- bb-mcp: build stage ---
FROM build-base AS build-bb-mcp

RUN CGO_ENABLED=0 go build -o /out/bb-mcp ./cmd/bb-mcp

# --- bb-mcp: minimal runtime ---
FROM scratch AS bb-mcp

COPY --from=build-bb-mcp /out/bb-mcp /bb-mcp
COPY --from=build-bb-mcp /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER 65532:65532

LABEL io.modelcontextprotocol.server.name="io.github.FabianSchurig/bitbucket-mcp"

EXPOSE 8080
ENTRYPOINT ["/bb-mcp"]
