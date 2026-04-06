// Package config loads optional MCP server configuration from mcp_config.yaml.
//
// The configuration file controls which tools are exposed to clients at runtime,
// allowing operators to filter dangerous endpoints (e.g. DELETE) and override
// tool descriptions without recompiling.
//
// If no configuration file is found, the server runs with all tools enabled.
package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// DefaultConfigFile is the default filename looked for in the working directory.
const DefaultConfigFile = "mcp_config.yaml"

// Config holds the runtime MCP server configuration.
type Config struct {
	Server         ServerConfig               `yaml:"server"`
	ToolOverrides  map[string]ToolOverrideItem `yaml:"tool_overrides"`
}

// ServerConfig controls which tools are exposed.
type ServerConfig struct {
	AllowedMethods []string `yaml:"allowed_methods"`
	IgnoredTools   []string `yaml:"ignored_tools"`
}

// ToolOverrideItem holds per-tool overrides.
type ToolOverrideItem struct {
	Description string `yaml:"description"`
}

// Load reads the configuration from the given file path.
// If the file does not exist, it returns a default (permissive) configuration.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}
	return Parse(data)
}

// Parse parses a YAML config from raw bytes.
func Parse(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if cfg.ToolOverrides == nil {
		cfg.ToolOverrides = map[string]ToolOverrideItem{}
	}
	return &cfg, nil
}

// DefaultConfig returns a configuration that permits all methods and tools.
func DefaultConfig() *Config {
	return &Config{
		ToolOverrides: map[string]ToolOverrideItem{},
	}
}

// IsMethodAllowed reports whether the given HTTP method is permitted.
// If no allowed_methods are configured, all methods are allowed.
func (c *Config) IsMethodAllowed(method string) bool {
	if len(c.Server.AllowedMethods) == 0 {
		return true
	}
	upper := strings.ToUpper(method)
	for _, m := range c.Server.AllowedMethods {
		if strings.ToUpper(m) == upper {
			return true
		}
	}
	return false
}

// IsToolIgnored reports whether the given tool name is in the ignored list.
func (c *Config) IsToolIgnored(name string) bool {
	for _, n := range c.Server.IgnoredTools {
		if n == name {
			return true
		}
	}
	return false
}
