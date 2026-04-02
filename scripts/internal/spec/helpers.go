package spec

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var (
	CamelUpperBoundary = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	CamelUpperRun      = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	NonAlpha           = regexp.MustCompile(`[^a-zA-Z0-9]+`)
)

// ToCamel converts a hyphen/underscore-separated string to CamelCase.
func ToCamel(s string) string {
	parts := NonAlpha.Split(s, -1)
	var b strings.Builder
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		b.WriteString(string(runes))
	}
	return b.String()
}

// ReservedGoNames contains identifiers that would shadow Go packages or
// builtins used by generated code.
var ReservedGoNames = map[string]bool{
	"context": true, "fmt": true, "json": true, "strconv": true,
	"client": true, "handlers": true, "output": true, "cobra": true,
	"error": true, "string": true, "int": true, "bool": true,
	"cmd": true, "args": true, "err": true, "body": true, "all": true,
}

// ToGoName converts a parameter name to a Go variable name (lowerCamelCase).
func ToGoName(s string) string {
	cc := ToCamel(s)
	if len(cc) == 0 {
		return cc
	}
	runes := []rune(cc)
	runes[0] = unicode.ToLower(runes[0])
	name := string(runes)
	if ReservedGoNames[name] {
		return name + "Param"
	}
	return name
}

// FlagName converts an API parameter name to a CLI flag name.
func FlagName(s string) string {
	return strings.ReplaceAll(s, "_", "-")
}

// GoType converts an OpenAPI type to a Go type string.
func GoType(t string) string {
	switch t {
	case "integer":
		return "int"
	case "boolean":
		return "bool"
	default:
		return "string"
	}
}

// DefaultValue returns the Go zero-value literal for a Go type.
func DefaultValue(t string) string {
	switch t {
	case "int":
		return "0"
	case "bool":
		return "false"
	default:
		return `""`
	}
}

// GoStringLit returns a Go string literal, preferring backtick quoting.
func GoStringLit(s string) string {
	if !strings.Contains(s, "`") {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

// IsPaginated checks if an operation returns a paginated response.
func IsPaginated(op *Op) bool {
	for _, code := range []string{"200", "201"} {
		resp, ok := op.Responses[code]
		if !ok {
			continue
		}
		for _, mt := range resp.Content {
			ref := mt.Schema.Ref
			if strings.Contains(ref, "paginated_") || strings.HasSuffix(ref, "search_result_page") {
				return true
			}
		}
	}
	return false
}

// ToKebab converts a camelCase operationID to a kebab-case CLI command name.
func ToKebab(operationID string) string {
	kebab := CamelUpperBoundary.ReplaceAllString(operationID, "${1}-${2}")
	return strings.ToLower(CamelUpperRun.ReplaceAllString(kebab, "${1}-${2}"))
}

// ResolveBodyRef extracts the $ref from a request body, if present.
func ResolveBodyRef(rb *RequestBody) string {
	if rb == nil {
		return ""
	}
	for _, mt := range rb.Content {
		if mt.Schema.Ref != "" {
			return mt.Schema.Ref
		}
	}
	return ""
}

// ParamDef describes a single API parameter for code generation.
type ParamDef struct {
	Name     string // API parameter name (e.g., "workspace")
	In       string // "path" or "query"
	Type     string // OpenAPI type: "string", "integer", "boolean"
	Required bool
	Desc     string // optional description override
}

// ParamsToParamDefs converts OpenAPI parameters to ParamDef slice.
func ParamsToParamDefs(params []Parameter) []ParamDef {
	defs := make([]ParamDef, 0, len(params))
	for _, p := range params {
		defs = append(defs, ParamDef{
			Name:     p.Name,
			In:       p.In,
			Type:     p.Schema.Type,
			Required: p.Required && p.In == "path",
		})
	}
	return defs
}

// InjectPaginationParams adds page/pagelen parameters if not already present.
func InjectPaginationParams(params []ParamDef) []ParamDef {
	hasParam := func(name string) bool {
		for _, p := range params {
			if p.Name == name {
				return true
			}
		}
		return false
	}
	// Pre-allocate capacity for up to 2 additional params.
	result := make([]ParamDef, len(params), len(params)+2)
	copy(result, params)
	if !hasParam("page") {
		result = append(result, ParamDef{
			Name:     "page",
			In:       "query",
			Type:     "integer",
			Desc:     "Page number (query parameter)",
		})
	}
	if !hasParam("pagelen") {
		result = append(result, ParamDef{
			Name:     "pagelen",
			In:       "query",
			Type:     "integer",
			Desc:     "Number of items per page (query parameter)",
		})
	}
	return result
}

// OperationDef holds metadata for a single API operation.
// This is the shared intermediate representation consumed by both
// CLI command generators and MCP/Terraform tool generators.
type OperationDef struct {
	OperationID string
	Method      string
	Path        string
	Summary     string
	Description string
	Params      []ParamDef
	BodyFields  []BodyField
	HasBody     bool
	Paginated   bool
}

// BuildOperation creates an OperationDef from a path entry, method/op, and schema.
func BuildOperation(pe PathEntry, entry MethodOp, schema *Schema) OperationDef {
	op := entry.Op
	allParams := MergeParams(pe.PathItem.Parameters, op.Parameters)
	paramDefs := ParamsToParamDefs(allParams)

	var bodyFields []BodyField
	bodyRef := ResolveBodyRef(op.RequestBody)
	if bodyRef != "" && schema.Components.Schemas != nil {
		visited := make(map[string]bool)
		bodyFields = ResolveBodyFields(schema.Components.Schemas, bodyRef, "", visited)
		sort.Slice(bodyFields, func(i, j int) bool {
			return bodyFields[i].Path < bodyFields[j].Path
		})
	}

	paginated := IsPaginated(op)
	if paginated {
		paramDefs = InjectPaginationParams(paramDefs)
	}

	return OperationDef{
		OperationID: op.OperationID,
		Method:      entry.Method,
		Path:        pe.Path,
		Summary:     op.Summary,
		Description: op.Description,
		Params:      paramDefs,
		BodyFields:  bodyFields,
		HasBody:     op.RequestBody != nil,
		Paginated:   paginated,
	}
}

// BuildOperations builds all OperationDefs from a schema.
func BuildOperations(schema *Schema) []OperationDef {
	var ops []OperationDef
	for _, pe := range SortedPathEntries(schema.Paths) {
		for _, entry := range MethodOps(pe.PathItem) {
			if entry.Op == nil || entry.Op.OperationID == "" {
				continue
			}
			ops = append(ops, BuildOperation(pe, entry, schema))
		}
	}
	return ops
}

// DescribeParam returns a human-readable description for a parameter.
func DescribeParam(p ParamDef) string {
	if p.Desc != "" {
		return p.Desc
	}
	return fmt.Sprintf("%s (%s parameter)", p.Name, p.In)
}
