package spec

import (
	"fmt"
	"strings"
)

// BodyField describes a flattened request body field for code generation.
type BodyField struct {
	Path     string // dot-separated path, e.g., "content.raw"
	FlagName string // CLI flag name, e.g., "content-raw"
	GoName   string // Go variable name, e.g., "bodyContentRaw"
	GoType   string // "string", "int", "bool"
	Default  string // Go zero-value literal
	Desc     string // human-readable description
}

// ─── Body field helpers ───────────────────────────────────────────────────────

// BodyFlagName converts a dot-separated path to a CLI flag name.
func BodyFlagName(path string) string {
	s := strings.ReplaceAll(path, ".", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

// BodyGoName converts a dot-separated path to a Go variable name.
func BodyGoName(path string) string {
	combined := strings.ReplaceAll(path, ".", "_")
	return "body" + ToCamel(combined)
}

// MakeBodyField creates a BodyField from a field path, OpenAPI type, and description.
func MakeBodyField(path, oaType, desc string) BodyField {
	gt := GoType(oaType)
	if desc == "" {
		desc = path
	}
	return BodyField{
		Path:     path,
		FlagName: BodyFlagName(path),
		GoName:   BodyGoName(path),
		GoType:   gt,
		Default:  DefaultValue(gt),
		Desc:     desc,
	}
}

// ─── Body schema resolution ──────────────────────────────────────────────────

const schemaRefPrefix = "#/components/schemas/"

// skipAllOfRefs lists schema names in allOf that should be skipped during
// body field resolution (e.g., the generic "object" base schema).
var skipAllOfRefs = map[string]bool{
	"object": true,
}

// skipPropNames lists property names that are auto-managed by the API and
// should not be exposed as writable body fields (e.g., timestamps, links,
// computed counts).
var skipPropNames = map[string]bool{
	"links": true, "user": true, "author": true,
	"created_on": true, "updated_on": true,
	"rendered": true, "closed_by": true,
	"id": true, "html": true, "deleted": true,
	"participants": true, "comment_count": true,
	"task_count": true, "merge_commit": true,
	"queued": true, "summary": true,
	"resolved_on": true, "resolved_by": true,
}

// skipPropertyRefs lists schema reference names whose nested properties should
// not be inlined into body fields (complex linked entities like users,
// repositories, commits).
var skipPropertyRefs = map[string]bool{
	"account": true, "user": true, "team": true,
	"repository": true, "link": true,
	"account_links": true, "team_links": true, "user_links": true,
	"comment_resolution": true, "commitstatus": true,
	"pullrequest": true, "base_commit": true, "commit": true,
}

// refIdOnlySchemas lists schema names where only the "id" sub-field should be
// exposed (rather than inlining the full schema), used for referenced entities.
var refIdOnlySchemas = map[string]bool{
	"comment": true,
}

// ResolveBodyFields recursively resolves a $ref to a list of flattened body fields.
func ResolveBodyFields(schemas map[string]any, ref, prefix string, visited map[string]bool) []BodyField {
	name := strings.TrimPrefix(ref, schemaRefPrefix)
	if visited[name] {
		return nil
	}
	visited[name] = true
	defer func() { delete(visited, name) }()

	raw, ok := schemas[name]
	if !ok {
		return nil
	}
	schema, ok := raw.(map[string]any)
	if !ok {
		return nil
	}
	return resolveSchemaObj(schemas, schema, prefix, visited)
}

func resolveSchemaObj(schemas map[string]any, schema map[string]any, prefix string, visited map[string]bool) []BodyField {
	if allOfRaw, ok := schema["allOf"]; ok {
		return resolveAllOf(schemas, allOfRaw, prefix, visited)
	}

	propsRaw, ok := schema["properties"]
	if !ok {
		return nil
	}
	props, ok := propsRaw.(map[string]any)
	if !ok {
		return nil
	}
	return flattenProperties(schemas, props, prefix, visited)
}

func resolveAllOf(schemas map[string]any, allOfRaw any, prefix string, visited map[string]bool) []BodyField {
	allOf, _ := allOfRaw.([]any)
	var fields []BodyField
	for _, entry := range allOf {
		m, ok := entry.(map[string]any)
		if !ok {
			continue
		}
		ref, ok := m["$ref"].(string)
		if !ok {
			fields = append(fields, resolveSchemaObj(schemas, m, prefix, visited)...)
			continue
		}
		refName := strings.TrimPrefix(ref, schemaRefPrefix)
		if !skipAllOfRefs[refName] {
			fields = append(fields, ResolveBodyFields(schemas, ref, prefix, visited)...)
		}
	}
	return fields
}

func flattenProperties(schemas map[string]any, props map[string]any, prefix string, visited map[string]bool) []BodyField {
	var fields []BodyField
	for name, propRaw := range props {
		if skipPropNames[name] {
			continue
		}
		prop, ok := propRaw.(map[string]any)
		if !ok {
			continue
		}
		path := name
		if prefix != "" {
			path = prefix + "." + name
		}
		fields = append(fields, flattenProperty(schemas, name, path, prop, visited)...)
	}
	return fields
}

func flattenProperty(schemas map[string]any, name, path string, prop map[string]any, visited map[string]bool) []BodyField {
	if ref, ok := prop["$ref"].(string); ok {
		return resolveRefProperty(schemas, name, path, ref, visited)
	}
	if _, ok := prop["allOf"]; ok {
		return resolveSchemaObj(schemas, prop, path, visited)
	}

	typ, _ := prop["type"].(string)
	desc, _ := prop["description"].(string)
	desc = appendEnumValues(prop, desc)

	switch typ {
	case "string", "integer", "boolean":
		return []BodyField{MakeBodyField(path, typ, desc)}
	case "object":
		if subProps, ok := prop["properties"].(map[string]any); ok {
			return flattenProperties(schemas, subProps, path, visited)
		}
	}
	return nil
}

func resolveRefProperty(schemas map[string]any, name, path, ref string, visited map[string]bool) []BodyField {
	refName := strings.TrimPrefix(ref, schemaRefPrefix)
	if skipPropertyRefs[refName] {
		return nil
	}
	if refIdOnlySchemas[refName] || visited[refName] {
		return []BodyField{MakeBodyField(path+".id", "integer", fmt.Sprintf("ID of referenced %s", name))}
	}
	return ResolveBodyFields(schemas, ref, path, visited)
}

func appendEnumValues(prop map[string]any, desc string) string {
	enumRaw, ok := prop["enum"]
	if !ok {
		return desc
	}
	enumArr, ok := enumRaw.([]any)
	if !ok {
		return desc
	}
	vals := make([]string, 0, len(enumArr))
	for _, v := range enumArr {
		vals = append(vals, fmt.Sprintf("%v", v))
	}
	if desc != "" {
		desc += " "
	}
	return desc + "[" + strings.Join(vals, ", ") + "]"
}
