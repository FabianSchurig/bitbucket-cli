// Package spec provides shared OpenAPI schema types and parsing utilities
// used by the code generation pipeline (CLI commands, MCP tools, etc.).
package spec

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// ─── OpenAPI schema types ─────────────────────────────────────────────────────

// Schema represents a self-contained OpenAPI 3.0 schema file.
type Schema struct {
	OpenAPI    string              `yaml:"openapi"`
	Info       map[string]any      `yaml:"info"`
	Paths      map[string]PathItem `yaml:"paths"`
	Components ComponentsSection   `yaml:"components"`
}

// ComponentsSection holds the components of an OpenAPI schema.
type ComponentsSection struct {
	Schemas map[string]any `yaml:"schemas"`
}

// PathItem represents an OpenAPI path item with its operations.
type PathItem struct {
	Parameters []Parameter `yaml:"parameters"`
	Get        *Op         `yaml:"get"`
	Post       *Op         `yaml:"post"`
	Put        *Op         `yaml:"put"`
	Patch      *Op         `yaml:"patch"`
	Delete     *Op         `yaml:"delete"`
}

// Op represents a single OpenAPI operation (e.g., GET /pullrequests).
type Op struct {
	OperationID  string             `yaml:"operationId"`
	Summary      string             `yaml:"summary"`
	Description  string             `yaml:"description"`
	Tags         []string           `yaml:"tags"`
	Parameters   []Parameter        `yaml:"parameters"`
	RequestBody  *RequestBody       `yaml:"requestBody"`
	Responses    Responses          `yaml:"responses"`
	OAuth2Scopes []OAuth2ScopeEntry `yaml:"x-atlassian-oauth2-scopes"`
}

// OAuth2ScopeEntry represents an x-atlassian-oauth2-scopes entry.
type OAuth2ScopeEntry struct {
	State  string   `yaml:"state"`
	Scheme string   `yaml:"scheme"`
	Scopes []string `yaml:"scopes"`
}

// Parameter represents an OpenAPI parameter (path, query, etc.).
type Parameter struct {
	Name     string          `yaml:"name"`
	In       string          `yaml:"in"`
	Required bool            `yaml:"required"`
	Schema   ParameterSchema `yaml:"schema"`
}

// ParameterSchema holds the type of a parameter.
type ParameterSchema struct {
	Type string `yaml:"type"`
}

// RequestBody represents an OpenAPI request body.
type RequestBody struct {
	Required bool                 `yaml:"required"`
	Content  map[string]MediaType `yaml:"content"`
}

// Responses maps HTTP status codes to response definitions.
type Responses map[string]ResponseDef

// ResponseDef represents a single HTTP response.
type ResponseDef struct {
	Content map[string]MediaType `yaml:"content"`
}

// MediaType represents a media type in an OpenAPI schema.
type MediaType struct {
	Schema RefSchema `yaml:"schema"`
}

// RefSchema holds a $ref pointer to a schema definition.
type RefSchema struct {
	Ref string `yaml:"$ref"`
}

// ─── Schema loading ───────────────────────────────────────────────────────────

// LoadSchema reads and parses an OpenAPI schema YAML file.
func LoadSchema(path string) (*Schema, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading schema: %w", err)
	}
	var schema Schema
	if err := yaml.Unmarshal(raw, &schema); err != nil {
		return nil, fmt.Errorf("parsing schema: %w", err)
	}
	return &schema, nil
}

// ─── Path and parameter helpers ───────────────────────────────────────────────

// PathEntry holds a path and its PathItem for sorted iteration.
type PathEntry struct {
	Path     string
	PathItem PathItem
}

// SortedPathEntries returns schema paths sorted alphabetically.
func SortedPathEntries(paths map[string]PathItem) []PathEntry {
	entries := make([]PathEntry, 0, len(paths))
	for p, pi := range paths {
		entries = append(entries, PathEntry{p, pi})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})
	return entries
}

// MethodOp pairs an HTTP method with its operation.
type MethodOp struct {
	Method string
	Op     *Op
}

// MethodOps returns all non-nil operations for a PathItem in standard order.
func MethodOps(pi PathItem) []MethodOp {
	return []MethodOp{
		{"GET", pi.Get},
		{"POST", pi.Post},
		{"PUT", pi.Put},
		{"PATCH", pi.Patch},
		{"DELETE", pi.Delete},
	}
}

// MergeParams combines path-level and operation-level parameters,
// with operation parameters taking precedence over path parameters.
func MergeParams(pathParams, opParams []Parameter) []Parameter {
	opParamNames := make(map[string]bool, len(opParams))
	for _, p := range opParams {
		opParamNames[p.Name] = true
	}
	var merged []Parameter
	for _, p := range pathParams {
		if !opParamNames[p.Name] {
			merged = append(merged, p)
		}
	}
	return append(merged, opParams...)
}

// CommandMeta extracts CLI parent-command metadata from the schema info section.
// It looks for x-cli-command-* extension fields and falls back to defaults.
func CommandMeta(info map[string]any) (name, use, short, long string) {
	name, _ = info["x-cli-command-name"].(string)
	use, _ = info["x-cli-command-use"].(string)
	short, _ = info["x-cli-command-short"].(string)
	long, _ = info["x-cli-command-long"].(string)
	// Defaults for backward compatibility with schemas lacking x-cli-command-* fields
	if name == "" {
		name = "PR"
	}
	if use == "" {
		use = "pr"
	}
	if short == "" {
		short = "Manage Bitbucket pull requests"
	}
	if long == "" {
		long = "Commands for listing, creating, reading, and merging Bitbucket pull requests."
	}
	return
}

// CommandCategory derives a human-readable category from the schema info title.
// It strips the "Bitbucket " prefix and " CLI" suffix (e.g., "Bitbucket Pull Requests CLI" → "Pull Requests").
func CommandCategory(info map[string]any) string {
	title, _ := info["title"].(string)
	title = strings.TrimPrefix(title, "Bitbucket ")
	title = strings.TrimSuffix(title, " CLI")
	if title == "" {
		return "Other"
	}
	return title
}
