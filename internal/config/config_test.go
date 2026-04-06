package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig_AllowsEverything(t *testing.T) {
	cfg := DefaultConfig()

	for _, method := range []string{"GET", "POST", "PUT", "PATCH", "DELETE"} {
		if !cfg.IsMethodAllowed(method) {
			t.Errorf("default config should allow %s", method)
		}
	}
	if cfg.IsToolIgnored("anything") {
		t.Error("default config should not ignore any tool")
	}
}

func TestLoad_MissingFile_ReturnsDefault(t *testing.T) {
	cfg, err := Load("/nonexistent/mcp_config.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.IsMethodAllowed("DELETE") {
		t.Error("missing file should return permissive default")
	}
}

func TestLoad_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mcp_config.yaml")
	content := `
server:
  allowed_methods: ["GET", "POST"]
  ignored_tools:
    - delete_repository

tool_overrides:
  get_pullrequests:
    description: "Custom description"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cfg.IsMethodAllowed("GET") {
		t.Error("GET should be allowed")
	}
	if !cfg.IsMethodAllowed("get") {
		t.Error("case-insensitive GET should be allowed")
	}
	if cfg.IsMethodAllowed("DELETE") {
		t.Error("DELETE should not be allowed")
	}
	if !cfg.IsToolIgnored("delete_repository") {
		t.Error("delete_repository should be ignored")
	}
	if cfg.IsToolIgnored("get_pullrequests") {
		t.Error("get_pullrequests should not be ignored")
	}
	override, ok := cfg.ToolOverrides["get_pullrequests"]
	if !ok {
		t.Fatal("expected tool override for get_pullrequests")
	}
	if override.Description != "Custom description" {
		t.Errorf("expected 'Custom description', got %q", override.Description)
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	_, err := Parse([]byte(`{invalid`))
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestParse_EmptyOverrides(t *testing.T) {
	cfg, err := Parse([]byte(`
server:
  allowed_methods: ["GET"]
`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ToolOverrides == nil {
		t.Error("ToolOverrides should be initialized to empty map")
	}
}
