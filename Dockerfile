# ============================================================================
# Dockerfile for bitbucket-cli
#
# Uses the hardened dhi.io/golang base image and installs pre-built binaries
# from the Go module proxy — no local source build required.
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
#   To add a new binary target, add a new stage that installs the binary
#   with `go install` (use the existing stages as a template), then add
#   the target to the build matrix in .github/workflows/docker.yml.
# ============================================================================

# --- bb-cli: hardened image for the Bitbucket CLI ---
FROM dhi.io/golang:1 AS bb-cli

RUN go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-cli@latest

ENTRYPOINT ["bb-cli"]

# --- bb-mcp: hardened image for the Bitbucket MCP server ---
FROM dhi.io/golang:1 AS bb-mcp

RUN go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-mcp@latest

EXPOSE 8080
ENTRYPOINT ["bb-mcp"]
