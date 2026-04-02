# ============================================================================
# Multi-stage Dockerfile for bitbucket-cli
#
# Targets:
#   bb-cli  — Bitbucket CLI (default)
#   bb-mcp  — Bitbucket MCP server
#
# Build examples:
#   docker build --target bb-cli -t bb-cli .
#   docker build --target bb-mcp -t bb-mcp .
#
# Extending this Dockerfile:
#   To add a new binary target:
#   1. Add a `RUN CGO_ENABLED=0 go build ...` line in the "builder" stage
#      that compiles the new binary into /out/.
#   2. Add a new final stage (use the bb-cli/bb-mcp stages as a template)
#      that copies the binary from the builder and sets the ENTRYPOINT.
#   3. Update .github/workflows/docker.yml to include the new target in the
#      build matrix so CI builds and pushes the image automatically.
# ============================================================================

# --- Builder stage: compile all binaries ---
FROM dhi.io/golang:1 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd/ cmd/
COPY internal/ internal/

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/bb-cli  ./cmd/bb-cli && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/bb-mcp  ./cmd/bb-mcp

# --- bb-cli: minimal hardened image for the CLI ---
FROM gcr.io/distroless/static-debian12:nonroot AS bb-cli

COPY --from=builder /out/bb-cli /usr/local/bin/bb-cli

ENTRYPOINT ["bb-cli"]

# --- bb-mcp: minimal hardened image for the MCP server ---
FROM gcr.io/distroless/static-debian12:nonroot AS bb-mcp

COPY --from=builder /out/bb-mcp /usr/local/bin/bb-mcp

EXPOSE 8080
ENTRYPOINT ["bb-mcp"]
