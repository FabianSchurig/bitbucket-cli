.PHONY: build test test-python lint fmt vet check generate generate-docs clean install

# Build all binaries
build:
	go build ./...

# Install all binaries to $GOPATH/bin (or $HOME/go/bin)
install:
	go install ./cmd/bb-cli
	go install ./cmd/bb-mcp

# Run Go tests with race detector and Python unit tests
test:
	go test -race ./...
	python3 -m unittest discover -s scripts -p '*_test.py'

# Run Python unit tests
test-python:
	python3 -m unittest discover -s scripts -p '*_test.py'

# Run golangci-lint (install via: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	golangci-lint run ./...

# Format all Go files
fmt:
	gofmt -w .

# Run go vet
vet:
	go vet ./...

# Check formatting without modifying files (CI-friendly)
fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "Files not formatted:"; gofmt -l .; exit 1)

# Run all checks (build + vet + lint + fmt-check + test with race)
check: build vet lint fmt-check test

# Regenerate all code from schemas (requires schemas in schema/)
generate:
	@if ! ls schema/*-schema.yaml >/dev/null 2>&1; then \
		echo "Error: no schema files found in schema/"; exit 1; \
	fi
	@for schema in schema/*-schema.yaml; do \
		base=$$(basename "$$schema" -schema.yaml); \
		if [ "$$base" = "pr" ]; then \
			go run scripts/gen_commands/main.go "$$schema" internal/commands/commands.gen.go; \
			go run scripts/gen_mcptools/main.go "$$schema" internal/mcptools/pr.gen.go; \
			go run scripts/gen_terraform/main.go "$$schema" internal/tfprovider/pr.gen.go; \
		else \
			go run scripts/gen_commands/main.go "$$schema" "internal/commands/$${base}.gen.go"; \
			go run scripts/gen_mcptools/main.go "$$schema" "internal/mcptools/$${base}.gen.go"; \
			go run scripts/gen_terraform/main.go "$$schema" "internal/tfprovider/$${base}.gen.go"; \
		fi; \
	done
	@echo "All code regenerated"

# Clean build artifacts
clean:
	rm -f bb-cli bb-mcp bb-tf gen_commands gen_mcptools gen_terraform main
	rm -rf dist/

# Generate Terraform provider docs, examples, and test files
# and the MCP tool reference page.
generate-docs:
	go run scripts/gen_tfdocs/main.go
	go run scripts/gen_mcpdocs/main.go
