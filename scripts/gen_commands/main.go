// gen_commands reads a PR schema YAML file and generates Cobra command boilerplate
// for each operation, writing the output to the specified Go file.
//
// Usage: go run scripts/gen_commands/main.go <schema.yaml> <output.go>
//
// The generated file defines NewPRCommand() which wires all operations as
// sub-commands under a "pr" parent command. Each generated command:
//   - Reads path parameters as required --flags
//   - Reads query parameters as optional --flags (typed from schema)
//   - For POST/PUT/PATCH operations, generates typed --flags from the request body schema
//   - Falls back to --body for raw JSON input (advanced)
//   - Calls handlers.Dispatch with method, URL template, params, etc.
//
// This script is run in CI whenever the schema changes.
package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"gopkg.in/yaml.v3"
)

var (
	camelUpperBoundary = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	camelUpperRun      = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
)

// ─── Schema types ─────────────────────────────────────────────────────────────

type Schema struct {
	OpenAPI    string              `yaml:"openapi"`
	Info       map[string]any      `yaml:"info"`
	Paths      map[string]PathItem `yaml:"paths"`
	Components ComponentsSection   `yaml:"components"`
}

type ComponentsSection struct {
	Schemas map[string]any `yaml:"schemas"`
}

type PathItem struct {
	Parameters []Parameter `yaml:"parameters"`
	Get        *Op         `yaml:"get"`
	Post       *Op         `yaml:"post"`
	Put        *Op         `yaml:"put"`
	Patch      *Op         `yaml:"patch"`
	Delete     *Op         `yaml:"delete"`
}

type Op struct {
	OperationID string       `yaml:"operationId"`
	Summary     string       `yaml:"summary"`
	Description string       `yaml:"description"`
	Tags        []string     `yaml:"tags"`
	Parameters  []Parameter  `yaml:"parameters"`
	RequestBody *RequestBody `yaml:"requestBody"`
	Responses   Responses    `yaml:"responses"`
}

type Parameter struct {
	Name     string          `yaml:"name"`
	In       string          `yaml:"in"`
	Required bool            `yaml:"required"`
	Schema   ParameterSchema `yaml:"schema"`
}

type ParameterSchema struct {
	Type string `yaml:"type"`
}

type RequestBody struct {
	Required bool                 `yaml:"required"`
	Content  map[string]MediaType `yaml:"content"`
}

type Responses map[string]ResponseDef

type ResponseDef struct {
	Content map[string]MediaType `yaml:"content"`
}

type MediaType struct {
	Schema RefSchema `yaml:"schema"`
}

type RefSchema struct {
	Ref string `yaml:"$ref"`
}

// ─── Template data types ──────────────────────────────────────────────────────

type CommandData struct {
	OperationID string
	Use         string
	Short       string
	Long        string
	Method      string
	Path        string
	Flags       []FlagData
	BodyFields  []BodyField
	HasBody     bool
	Paginated   bool
}

type FlagData struct {
	Name     string
	GoName   string
	GoType   string
	Default  string
	Usage    string
	Required bool
	In       string
	RawName  string
}

type BodyField struct {
	Path     string // "content.raw"
	FlagName string // "content-raw"
	GoName   string // "bodyContentRaw"
	GoType   string // "string", "int", "bool"
	Default  string // `""`, "0", "false"
	Desc     string
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

var nonAlpha = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func toCamel(s string) string {
	parts := nonAlpha.Split(s, -1)
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

func toGoName(s string) string {
	cc := toCamel(s)
	if len(cc) == 0 {
		return cc
	}
	runes := []rune(cc)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func flagName(s string) string {
	return strings.ReplaceAll(s, "_", "-")
}

func goType(t string) string {
	switch t {
	case "integer":
		return "int"
	case "boolean":
		return "bool"
	default:
		return "string"
	}
}

func defaultValue(t string) string {
	switch t {
	case "int":
		return "0"
	case "bool":
		return "false"
	default:
		return `""`
	}
}

func goStringLit(s string) string {
	if !strings.Contains(s, "`") {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

func isPaginated(op *Op) bool {
	for _, code := range []string{"200", "201"} {
		resp, ok := op.Responses[code]
		if !ok {
			continue
		}
		for _, mt := range resp.Content {
			if strings.Contains(mt.Schema.Ref, "paginated_") {
				return true
			}
		}
	}
	return false
}

// ─── Body field helpers ───────────────────────────────────────────────────────

func bodyFlagName(path string) string {
	s := strings.ReplaceAll(path, ".", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

func bodyGoName(path string) string {
	combined := strings.ReplaceAll(path, ".", "_")
	return "body" + toCamel(combined)
}

func makeBodyField(path, oaType, desc string) BodyField {
	gt := goType(oaType)
	if desc == "" {
		desc = path
	}
	return BodyField{
		Path:     path,
		FlagName: bodyFlagName(path),
		GoName:   bodyGoName(path),
		GoType:   gt,
		Default:  defaultValue(gt),
		Desc:     desc,
	}
}

// ─── Body schema resolution ──────────────────────────────────────────────────

const schemaRefPrefix = "#/components/schemas/"

var skipAllOfRefs = map[string]bool{
	"object": true,
}

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

var skipPropertyRefs = map[string]bool{
	"account": true, "user": true, "team": true,
	"repository": true, "link": true,
	"account_links": true, "team_links": true, "user_links": true,
	"comment_resolution": true, "commitstatus": true,
	"pullrequest": true, "base_commit": true, "commit": true,
}

var refIdOnlySchemas = map[string]bool{
	"comment": true,
}

func resolveBodyRef(rb *RequestBody) string {
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

func resolveBodyFields(schemas map[string]any, ref, prefix string, visited map[string]bool) []BodyField {
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
			fields = append(fields, resolveBodyFields(schemas, ref, prefix, visited)...)
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
		return []BodyField{makeBodyField(path, typ, desc)}
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
		return []BodyField{makeBodyField(path+".id", "integer", fmt.Sprintf("ID of referenced %s", name))}
	}
	return resolveBodyFields(schemas, ref, path, visited)
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

// ─── Code generation template ─────────────────────────────────────────────────

const fileTemplate = `// Code generated by scripts/gen_commands/main.go DO NOT EDIT.
// Source: {{.SchemaPath}}
//
// This file is regenerated whenever the schema changes in CI.
// Do not edit manually — run: go run scripts/gen_commands/main.go {{.SchemaPath}} <output.go>

package commands

import (
"context"
"encoding/json"
"fmt"
"strconv"

"github.com/spf13/cobra"

"github.com/FabianSchurig/bitbucket-cli/internal/client"
"github.com/FabianSchurig/bitbucket-cli/internal/handlers"
"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

// Ensure imports are used.
var (
_ = context.Background
_ = fmt.Errorf
_ = json.Marshal
_ = strconv.Itoa
_ = client.NewClient
_ = handlers.Dispatch
_ = output.Format
)

// New{{.CommandName}}Command returns the "{{.CommandUse}}" cobra command with all sub-commands registered.
func New{{.CommandName}}Command() *cobra.Command {
cmd := &cobra.Command{
Use:   "{{.CommandUse}}",
Short: {{goStringLit .CommandShort}},
Long:  {{goStringLit .CommandLong}},
}

cmd.AddCommand(
{{- range .Commands}}
new{{toCamel .OperationID}}Cmd(),
{{- end}}
)

return cmd
}
{{range .Commands}}
// new{{toCamel .OperationID}}Cmd returns the "{{$.CommandUse}} {{.Use}}" cobra command.
// operationId: {{.OperationID}}
func new{{toCamel .OperationID}}Cmd() *cobra.Command {
var (
{{- range .Flags}}
{{.GoName}} {{.GoType}}
{{- end}}
{{- range .BodyFields}}
{{.GoName}} {{.GoType}}
{{- end}}
{{- if .HasBody}}
body string
{{- end}}
{{- if .Paginated}}
all bool
{{- end}}
)

cmd := &cobra.Command{
Use:   "{{.Use}}",
Short: {{goStringLit .Short}},
Long:  {{goStringLit .Long}},
RunE: func(cmd *cobra.Command, args []string) error {
{{- range .Flags}}
{{- if .Required}}
{{- if eq .GoType "int"}}
if {{.GoName}} == 0 {
return fmt.Errorf("--{{.Name}} is required")
}
{{- else if eq .GoType "string"}}
if {{.GoName}} == "" {
return fmt.Errorf("--{{.Name}} is required")
}
{{- end}}
{{- end}}
{{- end}}
c, err := client.NewClient()
if err != nil {
return err
}
pathParams := map[string]string{
{{- range .Flags}}
{{- if eq .In "path"}}
{{- if eq .GoType "int"}}
"{{.RawName}}": strconv.Itoa({{.GoName}}),
{{- else}}
"{{.RawName}}": {{.GoName}},
{{- end}}
{{- end}}
{{- end}}
}
queryParams := map[string]string{
{{- range .Flags}}
{{- if eq .In "query"}}
{{- if eq .GoType "int"}}
"{{.RawName}}": strconv.Itoa({{.GoName}}),
{{- else if eq .GoType "bool"}}
"{{.RawName}}": strconv.FormatBool({{.GoName}}),
{{- else}}
"{{.RawName}}": {{.GoName}},
{{- end}}
{{- end}}
{{- end}}
}
{{- if .HasBody}}
if body == "" {
bodyObj := map[string]any{}
{{- range .BodyFields}}
{{- if eq .GoType "string"}}
if {{.GoName}} != "" {
handlers.SetNested(bodyObj, "{{.Path}}", {{.GoName}})
}
{{- else if eq .GoType "int"}}
if {{.GoName}} != 0 {
handlers.SetNested(bodyObj, "{{.Path}}", {{.GoName}})
}
{{- else if eq .GoType "bool"}}
if {{.GoName}} {
handlers.SetNested(bodyObj, "{{.Path}}", {{.GoName}})
}
{{- end}}
{{- end}}
if len(bodyObj) > 0 {
b, _ := json.Marshal(bodyObj)
body = string(b)
}
}
{{- end}}
{{- if not .HasBody}}
body := ""
{{- end}}
return handlers.Dispatch(context.Background(), c, handlers.Request{
					Method:      "{{.Method}}",
					URLTemplate: "{{.Path}}",
					PathParams:  pathParams,
					QueryParams: queryParams,
					Body:        body,
					All:         {{if .Paginated}}all{{else}}false{{end}},
				})
},
}

{{- range .Flags}}
cmd.Flags().{{flagFuncName .GoType}}Var(&{{.GoName}}, "{{.Name}}", {{.Default}}, "{{.Usage}}")
{{- end}}
{{- range .BodyFields}}
cmd.Flags().{{flagFuncName .GoType}}Var(&{{.GoName}}, "{{.FlagName}}", {{.Default}}, {{goStringLit .Desc}})
{{- end}}
{{- if .HasBody}}
cmd.Flags().StringVar(&body, "body", "", "Raw JSON request body (advanced)")
{{- end}}
{{- if .Paginated}}
cmd.Flags().BoolVar(&all, "all", true, "Traverse all pages (follows 'next' cursor)")
{{- end}}
return cmd
}
{{end}}
`

// ─── Main ─────────────────────────────────────────────────────────────────────

type FileData struct {
	SchemaPath   string
	CommandName  string // e.g., "PR", "Hooks" — used in function name NewXxxCommand()
	CommandUse   string // e.g., "pr", "hooks" — Cobra Use field
	CommandShort string // Cobra Short field
	CommandLong  string // Cobra Long field
	Commands     []CommandData
}

type pathEntry struct {
	path     string
	pathItem PathItem
}

type methodOp struct {
	method string
	op     *Op
}

func loadSchema(path string) (*Schema, error) {
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

func sortedPathEntries(paths map[string]PathItem) []pathEntry {
	entries := make([]pathEntry, 0, len(paths))
	for p, pi := range paths {
		entries = append(entries, pathEntry{p, pi})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].path < entries[j].path
	})
	return entries
}

func mergeParams(pathParams, opParams []Parameter) []Parameter {
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

func paramsToFlags(params []Parameter) []FlagData {
	flags := make([]FlagData, 0, len(params))
	for _, p := range params {
		gt := goType(p.Schema.Type)
		flags = append(flags, FlagData{
			Name:     flagName(p.Name),
			GoName:   toGoName(p.Name),
			GoType:   gt,
			Default:  defaultValue(gt),
			Usage:    fmt.Sprintf("%s (%s parameter)", p.Name, p.In),
			Required: p.Required && p.In == "path",
			In:       p.In,
			RawName:  p.Name,
		})
	}
	return flags
}

func injectPaginationFlags(flags []FlagData) []FlagData {
	hasFlag := func(name string) bool {
		for _, f := range flags {
			if f.RawName == name {
				return true
			}
		}
		return false
	}
	if !hasFlag("page") {
		flags = append(flags, FlagData{
			Name: "page", GoName: "page", GoType: "int",
			Default: "0", Usage: "Page number (query parameter)",
			In: "query", RawName: "page",
		})
	}
	if !hasFlag("pagelen") {
		flags = append(flags, FlagData{
			Name: "pagelen", GoName: "pagelen", GoType: "int",
			Default: "0", Usage: "Number of items per page (query parameter)",
			In: "query", RawName: "pagelen",
		})
	}
	return flags
}

func buildCommand(pe pathEntry, entry methodOp, schema *Schema) CommandData {
	op := entry.op
	allParams := mergeParams(pe.pathItem.Parameters, op.Parameters)
	flags := paramsToFlags(allParams)

	var bodyFields []BodyField
	bodyRef := resolveBodyRef(op.RequestBody)
	if bodyRef != "" && schema.Components.Schemas != nil {
		visited := make(map[string]bool)
		bodyFields = resolveBodyFields(schema.Components.Schemas, bodyRef, "", visited)
		sort.Slice(bodyFields, func(i, j int) bool {
			return bodyFields[i].Path < bodyFields[j].Path
		})
	}

	paginated := isPaginated(op)
	if paginated {
		flags = injectPaginationFlags(flags)
	}

	// Two-pass camelCase → kebab-case conversion.
	// Pass 1 (camelUpperBoundary): splits lower→upper boundaries (e.g. "getA" → "get-A").
	// Pass 2 (camelUpperRun): splits consecutive uppercase runs (e.g. "AWebhook" → "A-Webhook").
	// Together they handle single-letter words like "A" in "getAWebhookResource" → "get-a-webhook-resource".
	kebab := camelUpperBoundary.ReplaceAllString(op.OperationID, "${1}-${2}")
	kebab = strings.ToLower(camelUpperRun.ReplaceAllString(kebab, "${1}-${2}"))

	return CommandData{
		OperationID: op.OperationID,
		Use:         kebab,
		Short:       op.Summary,
		Long:        op.Description,
		Method:      entry.method,
		Path:        pe.path,
		Flags:       flags,
		BodyFields:  bodyFields,
		HasBody:     op.RequestBody != nil,
		Paginated:   paginated,
	}
}

func buildCommands(schema *Schema) []CommandData {
	var commands []CommandData
	for _, pe := range sortedPathEntries(schema.Paths) {
		pathItem := pe.pathItem
		for _, entry := range []methodOp{
			{"GET", pathItem.Get},
			{"POST", pathItem.Post},
			{"PUT", pathItem.Put},
			{"PATCH", pathItem.Patch},
			{"DELETE", pathItem.Delete},
		} {
			if entry.op == nil || entry.op.OperationID == "" {
				continue
			}
			commands = append(commands, buildCommand(pe, entry, schema))
		}
	}
	return commands
}

var templateFuncMap = template.FuncMap{
	"toCamel":     toCamel,
	"goStringLit": goStringLit,
	"flagFuncName": func(goType string) string {
		switch goType {
		case "int":
			return "Int"
		case "bool":
			return "Bool"
		default:
			return "String"
		}
	},
}

func generate(data FileData, outputPath string) error {
	tmpl, err := template.New("commands").Funcs(templateFuncMap).Parse(fileTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

// commandMeta extracts CLI parent-command metadata from the schema info section.
// It looks for x-cli-command-* extension fields and falls back to defaults.
func commandMeta(info map[string]any) (name, use, short, long string) {
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

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <schema.yaml> <output.go>\n", os.Args[0])
		os.Exit(1)
	}

	schemaPath := os.Args[1]
	outputPath := os.Args[2]

	schema, err := loadSchema(schemaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	cmdName, cmdUse, cmdShort, cmdLong := commandMeta(schema.Info)

	data := FileData{
		SchemaPath:   schemaPath,
		CommandName:  cmdName,
		CommandUse:   cmdUse,
		CommandShort: cmdShort,
		CommandLong:  cmdLong,
		Commands:     buildCommands(schema),
	}

	if err := generate(data, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %d commands → %s\n", len(data.Commands), outputPath)
}
