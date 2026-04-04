package spec

import (
	"fmt"
	"strings"
)

// BodyField describes a flattened request body field for code generation.
type BodyField struct {
	Path       string      // dot-separated path, e.g., "content.raw"
	FlagName   string      // CLI flag name, e.g., "content-raw"
	GoName     string      // Go variable name, e.g., "bodyContentRaw"
	GoType     string      // "string", "int", "bool"
	Default    string      // Go zero-value literal
	Desc       string      // human-readable description
	IsArray    bool        // true when the field is an array
	IsObject   bool        // true when the field is a nested object ($ref or inline)
	ItemFields []BodyField // nested fields for array item objects or object properties
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

// FieldResolveOpts controls which properties are skipped during field resolution.
type FieldResolveOpts struct {
	SkipPropNames map[string]bool
	SkipAllOfRefs map[string]bool
	// MaxRefDepth controls how many levels of property $ref references are
	// followed. allOf composition does not count as a level (it is structural
	// inheritance, not a relationship). Use 1 to expose one level of nested
	// entity properties (e.g., target.hash for a commit ref), 0 to not follow
	// any property $ref references at all.
	MaxRefDepth int
}

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

// responseSkipPropNames is more permissive than skipPropNames — response fields
// like created_on, updated_on, comment_count etc. are useful computed values.
var responseSkipPropNames = map[string]bool{
	"links": true, "html": true, "rendered": true,
}

// arrayItemSkipPropNames lists properties to skip when resolving array item fields.
// More permissive than body opts — includes identifiers like id, uuid — but still
// skips complex objects like links.
var arrayItemSkipPropNames = map[string]bool{
	"links": true, "html": true, "rendered": true,
}

// BodyFieldOpts returns the default options for request body field resolution.
// MaxRefDepth=1 follows one level of property $ref references, exposing
// direct properties of referenced entities (e.g., target.hash for commit).
func BodyFieldOpts() FieldResolveOpts {
	return FieldResolveOpts{
		SkipPropNames: skipPropNames,
		SkipAllOfRefs: skipAllOfRefs,
		MaxRefDepth:   1,
	}
}

// ResponseFieldOpts returns options for response field resolution (more permissive).
// MaxRefDepth=1 follows one level of property $ref references.
func ResponseFieldOpts() FieldResolveOpts {
	return FieldResolveOpts{
		SkipPropNames: responseSkipPropNames,
		SkipAllOfRefs: skipAllOfRefs,
		MaxRefDepth:   1,
	}
}

// ArrayItemFieldOpts returns options for resolving fields inside array item schemas.
// MaxRefDepth=0 keeps array items shallow — only direct properties are exposed.
func ArrayItemFieldOpts() FieldResolveOpts {
	return FieldResolveOpts{
		SkipPropNames: arrayItemSkipPropNames,
		SkipAllOfRefs: skipAllOfRefs,
		MaxRefDepth:   0,
	}
}

// ResolveBodyFields recursively resolves a $ref to a list of flattened body fields.
func ResolveBodyFields(schemas map[string]any, ref, prefix string, visited map[string]bool) []BodyField {
	return ResolveFields(schemas, ref, prefix, visited, BodyFieldOpts())
}

// ResolveResponseFields resolves a response schema $ref to a list of flattened fields
// using more permissive skip lists (e.g., includes created_on, updated_on, etc.).
func ResolveResponseFields(schemas map[string]any, ref, prefix string, visited map[string]bool) []BodyField {
	return ResolveFields(schemas, ref, prefix, visited, ResponseFieldOpts())
}

// FlattenBodyFields recursively flattens nested object BodyFields into flat
// dot-separated paths. Array fields are preserved as-is (only object nesting
// is flattened). This is used by CLI and MCP generators that need flat flags.
func FlattenBodyFields(fields []BodyField) []BodyField {
	var flat []BodyField
	for _, f := range fields {
		flat = append(flat, flattenBodyFieldRec(f, "")...)
	}
	return flat
}

func flattenBodyFieldRec(f BodyField, prefix string) []BodyField {
	path := f.Path
	if prefix != "" {
		path = prefix + "." + f.Path
	}

	if f.IsObject && len(f.ItemFields) > 0 {
		var flat []BodyField
		for _, item := range f.ItemFields {
			flat = append(flat, flattenBodyFieldRec(item, path)...)
		}
		return flat
	}

	// For scalars and arrays, reconstruct with full path.
	result := f
	result.Path = path
	result.FlagName = BodyFlagName(path)
	result.GoName = BodyGoName(path)
	return []BodyField{result}
}

// ResolveFields recursively resolves a $ref to a list of flattened fields
// with configurable skip lists via opts.
func ResolveFields(schemas map[string]any, ref, prefix string, visited map[string]bool, opts FieldResolveOpts) []BodyField {
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
	return resolveSchemaObj(schemas, schema, prefix, visited, opts)
}

func resolveSchemaObj(schemas map[string]any, schema map[string]any, prefix string, visited map[string]bool, opts FieldResolveOpts) []BodyField {
	if allOfRaw, ok := schema["allOf"]; ok {
		return resolveAllOf(schemas, allOfRaw, prefix, visited, opts)
	}

	propsRaw, ok := schema["properties"]
	if !ok {
		return nil
	}
	props, ok := propsRaw.(map[string]any)
	if !ok {
		return nil
	}
	return flattenProperties(schemas, props, prefix, visited, opts)
}

func resolveAllOf(schemas map[string]any, allOfRaw any, prefix string, visited map[string]bool, opts FieldResolveOpts) []BodyField {
	allOf, _ := allOfRaw.([]any)
	var fields []BodyField
	for _, entry := range allOf {
		m, ok := entry.(map[string]any)
		if !ok {
			continue
		}
		ref, ok := m["$ref"].(string)
		if !ok {
			fields = append(fields, resolveSchemaObj(schemas, m, prefix, visited, opts)...)
			continue
		}
		refName := strings.TrimPrefix(ref, schemaRefPrefix)
		if !opts.SkipAllOfRefs[refName] {
			fields = append(fields, ResolveFields(schemas, ref, prefix, visited, opts)...)
		}
	}
	return fields
}

func flattenProperties(schemas map[string]any, props map[string]any, prefix string, visited map[string]bool, opts FieldResolveOpts) []BodyField {
	var fields []BodyField
	for name, propRaw := range props {
		if opts.SkipPropNames[name] {
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
		fields = append(fields, flattenProperty(schemas, name, path, prop, visited, opts)...)
	}
	return fields
}

func flattenProperty(schemas map[string]any, name, path string, prop map[string]any, visited map[string]bool, opts FieldResolveOpts) []BodyField {
	if ref, ok := prop["$ref"].(string); ok {
		return resolveRefProperty(schemas, name, path, ref, prop, visited, opts)
	}
	if _, ok := prop["allOf"]; ok {
		nested := resolveSchemaObj(schemas, prop, "", visited, opts)
		if len(nested) > 0 {
			desc, _ := prop["description"].(string)
			if desc == "" {
				desc = name
			}
			return []BodyField{{
				Path: path, FlagName: BodyFlagName(path), GoName: BodyGoName(path),
				GoType: "string", Default: `""`, Desc: desc,
				IsObject: true, ItemFields: nested,
			}}
		}
		return nil
	}

	typ, _ := prop["type"].(string)
	desc, _ := prop["description"].(string)
	desc = appendEnumValues(prop, desc)

	switch typ {
	case "string", "integer", "boolean":
		return []BodyField{MakeBodyField(path, typ, desc)}
	case "object":
		if subProps, ok := prop["properties"].(map[string]any); ok {
			nested := flattenProperties(schemas, subProps, "", visited, opts)
			if len(nested) > 0 {
				if desc == "" {
					desc = name
				}
				return []BodyField{{
					Path: path, FlagName: BodyFlagName(path), GoName: BodyGoName(path),
					GoType: "string", Default: `""`, Desc: desc,
					IsObject: true, ItemFields: nested,
				}}
			}
		}
	case "array":
		items, _ := prop["items"].(map[string]any)
		if items != nil {
			itemOpts := ArrayItemFieldOpts()
			// Copy the parent visited map to prevent infinite recursion with
			// self-referencing schemas (e.g., base_commit.parents → base_commit).
			itemVisited := make(map[string]bool, len(visited))
			for k, v := range visited {
				itemVisited[k] = v
			}

			// Try to resolve array item fields from $ref.
			if ref, ok := items["$ref"].(string); ok {
				itemFields := ResolveFields(schemas, ref, "", itemVisited, itemOpts)
				if len(itemFields) > 0 {
					if desc == "" {
						desc = name
					}
					return []BodyField{{
						Path:       path,
						FlagName:   BodyFlagName(path),
						GoName:     BodyGoName(path),
						GoType:     "string",
						Default:    `""`,
						Desc:       desc,
						IsArray:    true,
						ItemFields: itemFields,
					}}
				}
			}

			// Try inline object definition: items: {type: object, properties: {...}}
			if itemType, _ := items["type"].(string); itemType == "object" {
				if subProps, ok := items["properties"].(map[string]any); ok {
					itemFields := flattenProperties(schemas, subProps, "", itemVisited, itemOpts)
					if len(itemFields) > 0 {
						if desc == "" {
							desc = name
						}
						return []BodyField{{
							Path:       path,
							FlagName:   BodyFlagName(path),
							GoName:     BodyGoName(path),
							GoType:     "string",
							Default:    `""`,
							Desc:       desc,
							IsArray:    true,
							ItemFields: itemFields,
						}}
					}
				}
				// Inline object with allOf.
				if _, ok := items["allOf"]; ok {
					itemFields := resolveSchemaObj(schemas, items, "", itemVisited, itemOpts)
					if len(itemFields) > 0 {
						if desc == "" {
							desc = name
						}
						return []BodyField{{
							Path:       path,
							FlagName:   BodyFlagName(path),
							GoName:     BodyGoName(path),
							GoType:     "string",
							Default:    `""`,
							Desc:       desc,
							IsArray:    true,
							ItemFields: itemFields,
						}}
					}
				}
			}

			// Simple type arrays (string, integer, boolean) — no nested fields.
			if itemType, _ := items["type"].(string); itemType == "string" || itemType == "integer" || itemType == "boolean" {
				if desc == "" {
					desc = name
				}
				itemDesc := appendEnumValues(items, "")
				if itemDesc != "" {
					desc = desc + " " + itemDesc
				}
				return []BodyField{{
					Path:     path,
					FlagName: BodyFlagName(path),
					GoName:   BodyGoName(path),
					GoType:   "string",
					Default:  `""`,
					Desc:     desc,
					IsArray:  true,
					// ItemFields is nil — signals a simple list (List of String).
				}}
			}
		}
		// Fallback: expose array as a single string field accepting a JSON array.
		if desc == "" {
			desc = name
		}
		return []BodyField{MakeBodyField(path, "string", desc+" (JSON array)")}
	}
	return nil
}

func resolveRefProperty(schemas map[string]any, name, path, ref string, prop map[string]any, visited map[string]bool, opts FieldResolveOpts) []BodyField {
	if opts.MaxRefDepth <= 0 {
		if idField, ok := resolveReferencedIDField(schemas, ref); ok {
			desc, _ := prop["description"].(string)
			if desc == "" {
				desc = name
			}
			return []BodyField{{
				Path: path, FlagName: BodyFlagName(path), GoName: BodyGoName(path),
				GoType: "string", Default: `""`, Desc: desc,
				IsObject: true, ItemFields: []BodyField{idField},
			}}
		}
		return nil
	}
	deeper := opts
	deeper.MaxRefDepth = opts.MaxRefDepth - 1
	// Resolve with empty prefix so item fields have relative paths.
	nested := ResolveFields(schemas, ref, "", visited, deeper)
	if idField, ok := resolveReferencedIDField(schemas, ref); ok && !hasBodyFieldPath(nested, idField.Path) {
		nested = append(nested, idField)
	}
	if len(nested) == 0 {
		return nil
	}
	desc, _ := prop["description"].(string)
	if desc == "" {
		// Try to get description from the referenced schema.
		refName := strings.TrimPrefix(ref, schemaRefPrefix)
		if raw, ok := schemas[refName]; ok {
			if s, ok := raw.(map[string]any); ok {
				desc, _ = s["description"].(string)
			}
		}
	}
	if desc == "" {
		desc = name
	}
	return []BodyField{{
		Path: path, FlagName: BodyFlagName(path), GoName: BodyGoName(path),
		GoType: "string", Default: `""`, Desc: desc,
		IsObject: true, ItemFields: nested,
	}}
}

func hasBodyFieldPath(fields []BodyField, path string) bool {
	for _, field := range fields {
		if field.Path == path {
			return true
		}
	}
	return false
}

func resolveReferencedIDField(schemas map[string]any, ref string) (BodyField, bool) {
	prop, ok := resolveSchemaProperty(schemas, ref, "id", map[string]bool{})
	if !ok {
		return BodyField{}, false
	}
	propType, _ := prop["type"].(string)
	if propType == "" {
		propType = "string"
	}
	desc, _ := prop["description"].(string)
	return MakeBodyField("id", propType, desc), true
}

func resolveSchemaProperty(schemas map[string]any, ref, propName string, seen map[string]bool) (map[string]any, bool) {
	name := strings.TrimPrefix(ref, schemaRefPrefix)
	if seen[name] {
		return nil, false
	}
	seen[name] = true
	defer delete(seen, name)

	raw, ok := schemas[name]
	if !ok {
		return nil, false
	}
	schema, ok := raw.(map[string]any)
	if !ok {
		return nil, false
	}
	return findSchemaProperty(schemas, schema, propName, seen)
}

func findSchemaProperty(schemas map[string]any, schema map[string]any, propName string, seen map[string]bool) (map[string]any, bool) {
	if propsRaw, ok := schema["properties"]; ok {
		if props, ok := propsRaw.(map[string]any); ok {
			if propRaw, ok := props[propName]; ok {
				if prop, ok := propRaw.(map[string]any); ok {
					return prop, true
				}
			}
		}
	}
	allOfRaw, ok := schema["allOf"]
	if !ok {
		return nil, false
	}
	allOf, _ := allOfRaw.([]any)
	for _, entry := range allOf {
		m, ok := entry.(map[string]any)
		if !ok {
			continue
		}
		if ref, ok := m["$ref"].(string); ok {
			if prop, found := resolveSchemaProperty(schemas, ref, propName, seen); found {
				return prop, true
			}
			continue
		}
		if prop, found := findSchemaProperty(schemas, m, propName, seen); found {
			return prop, true
		}
	}
	return nil, false
}

// ResolveResponseRef extracts the response entity schema $ref from an operation.
// For paginated responses, it digs into the values array items to get the
// underlying entity schema. Returns empty string if no response schema found.
func ResolveResponseRef(op *Op, schemas map[string]any) string {
	if op == nil {
		return ""
	}
	for _, code := range []string{"200", "201"} {
		resp, ok := op.Responses[code]
		if !ok {
			continue
		}
		for _, mt := range resp.Content {
			ref := mt.Schema.Ref
			if ref == "" {
				continue
			}
			// For paginated responses, extract the item schema from values array.
			refName := strings.TrimPrefix(ref, schemaRefPrefix)
			if strings.Contains(refName, "paginated_") || strings.HasSuffix(refName, "search_result_page") {
				if itemRef := extractPaginatedItemRef(schemas, refName); itemRef != "" {
					return itemRef
				}
			}
			return ref
		}
	}
	return ""
}

// extractPaginatedItemRef extracts the $ref from a paginated schema's
// values.items field, e.g., paginated_projects → project.
func extractPaginatedItemRef(schemas map[string]any, schemaName string) string {
	raw, ok := schemas[schemaName]
	if !ok {
		return ""
	}
	schema, ok := raw.(map[string]any)
	if !ok {
		return ""
	}
	props, ok := schema["properties"].(map[string]any)
	if !ok {
		return ""
	}
	values, ok := props["values"].(map[string]any)
	if !ok {
		return ""
	}
	items, ok := values["items"].(map[string]any)
	if !ok {
		return ""
	}
	ref, _ := items["$ref"].(string)
	return ref
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
