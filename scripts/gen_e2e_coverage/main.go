// gen_e2e_coverage generates docs/e2e-coverage.md, a list of Bitbucket API
// endpoints exposed by the Terraform provider and whether they are exercised
// by a real-API acceptance test (TestAccRealAPI_*) in
// internal/tfprovider/acceptance_test.go.
//
// Usage: go run scripts/gen_e2e_coverage/main.go
//
// The generator:
//  1. Imports tfprovider to enumerate every registered ResourceGroup and the
//     CRUD operations (method + path) it exposes.
//  2. Parses internal/tfprovider/acceptance_test.go and, for each function
//     whose name starts with TestAccRealAPI_, collects every "bitbucket_<x>"
//     identifier referenced inside its body.
//  3. Marks a group as covered when any test references its terraform type
//     name ("bitbucket_" + snake_case(TypeName)).
//
// Output is deterministic (groups sorted alphabetically, tests sorted) so the
// file can be safely auto-committed by CI.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"sort"
	"strings"

	// Importing the tfprovider package triggers the init() calls in every
	// generated *.gen.go and in sub_resources.go, populating the registry.
	"github.com/FabianSchurig/bitbucket-cli/internal/tfprovider"
)

const (
	acceptanceTestFile = "internal/tfprovider/acceptance_test.go"
	outputFile         = "docs/e2e-coverage.md"
	testPrefix         = "TestAccRealAPI_"
)

// bitbucketTypeRefRE matches Terraform type identifiers (bitbucket_<snake>)
// referenced inside HCL config strings in the acceptance tests.
var bitbucketTypeRefRE = regexp.MustCompile(`bitbucket_[a-z0-9_]+`)

// operation represents a single API endpoint exposed by a resource group.
type operation struct {
	CRUD   string // "Create", "Read", "Update", "Delete", "List"
	Method string
	Path   string
}

// groupCoverage is the per-group view rendered into the markdown file.
type groupCoverage struct {
	TFName      string // e.g., "bitbucket_repos"
	Category    string
	Description string
	Ops         []operation
	Tests       []string // sorted TestAccRealAPI_* names referencing this group
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gen_e2e_coverage:", err)
		os.Exit(1)
	}
}

func run() error {
	testRefs, err := parseTestReferences(acceptanceTestFile)
	if err != nil {
		return fmt.Errorf("parse %s: %w", acceptanceTestFile, err)
	}

	groups := collectGroups(testRefs)
	out := renderMarkdown(groups)

	if err := os.WriteFile(outputFile, []byte(out), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", outputFile, err)
	}
	fmt.Printf("wrote %s (%d groups, %d covered)\n", outputFile, len(groups), countCovered(groups))
	return nil
}

// parseTestReferences walks every top-level func declaration in the given Go
// source file and, for those named TestAccRealAPI_*, returns a map of
// terraform type name → sorted list of test names that reference it.
func parseTestReferences(path string) (map[string][]string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	// tfName → set of test names
	refs := map[string]map[string]struct{}{}

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		name := fn.Name.Name
		if !strings.HasPrefix(name, testPrefix) {
			continue
		}
		// Walk the body and collect every string literal; HCL configs live in
		// raw or interpreted string literals passed to Config fields.
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			for _, m := range bitbucketTypeRefRE.FindAllString(lit.Value, -1) {
				if refs[m] == nil {
					refs[m] = map[string]struct{}{}
				}
				refs[m][name] = struct{}{}
			}
			return true
		})
	}

	out := map[string][]string{}
	for tfName, set := range refs {
		names := make([]string, 0, len(set))
		for n := range set {
			names = append(names, n)
		}
		sort.Strings(names)
		out[tfName] = names
	}
	return out, nil
}

// collectGroups returns the registered groups sorted by terraform name,
// each annotated with the tests that reference it.
func collectGroups(testRefs map[string][]string) []groupCoverage {
	registered := tfprovider.RegisteredGroups()
	groups := make([]groupCoverage, 0, len(registered))
	for _, g := range registered {
		tfName := "bitbucket_" + strings.ReplaceAll(g.TypeName, "-", "_")
		ops := collectOps(g.Ops)
		groups = append(groups, groupCoverage{
			TFName:      tfName,
			Category:    g.Category,
			Description: g.Description,
			Ops:         ops,
			Tests:       testRefs[tfName],
		})
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].TFName < groups[j].TFName })
	return groups
}

// collectOps flattens CRUDOps into an ordered slice of operation entries.
func collectOps(ops tfprovider.CRUDOps) []operation {
	var out []operation
	add := func(crud string, op *tfprovider.OperationDef) {
		if op == nil {
			return
		}
		out = append(out, operation{CRUD: crud, Method: op.Method, Path: op.Path})
	}
	add("Create", ops.Create)
	add("Read", ops.Read)
	add("Update", ops.Update)
	add("Delete", ops.Delete)
	add("List", ops.List)
	return out
}

func countCovered(groups []groupCoverage) int {
	n := 0
	for _, g := range groups {
		if len(g.Tests) > 0 {
			n++
		}
	}
	return n
}

// renderMarkdown produces the e2e-coverage.md content.
func renderMarkdown(groups []groupCoverage) string {
	var b strings.Builder
	covered := countCovered(groups)
	total := len(groups)

	b.WriteString("# Real-API E2E Test Coverage\n\n")
	b.WriteString("<!-- DO NOT EDIT. This file is auto-generated by `make generate-docs` -->\n")
	b.WriteString("<!-- (`go run scripts/gen_e2e_coverage/main.go`). It is refreshed -->\n")
	b.WriteString("<!-- automatically on push to `main` by the CI `format` job. -->\n\n")
	b.WriteString("This page lists every Terraform resource group exposed by the provider ")
	b.WriteString("and whether it is exercised by a real-API acceptance test (functions ")
	b.WriteString("named `TestAccRealAPI_*` in [`internal/tfprovider/acceptance_test.go`]")
	b.WriteString("(../internal/tfprovider/acceptance_test.go).\n\n")
	b.WriteString("A group counts as covered when at least one `TestAccRealAPI_*` test ")
	b.WriteString("references its Terraform type name (`bitbucket_<group>`) inside the ")
	b.WriteString("test's HCL configuration. The endpoints listed under each group are ")
	b.WriteString("the CRUD operations the provider wires up for that group; running the ")
	b.WriteString("referenced test exercises some or all of them against the real ")
	b.WriteString("Bitbucket Cloud API.\n\n")
	fmt.Fprintf(&b, "**Coverage: %d / %d resource groups (%d%%).**\n\n",
		covered, total, percent(covered, total))
	b.WriteString("To add coverage for a missing group, add a new `TestAccRealAPI_*` ")
	b.WriteString("function in `acceptance_test.go` that uses the corresponding ")
	b.WriteString("`bitbucket_<group>` resource or data source, then run ")
	b.WriteString("`make generate-docs` to refresh this file.\n\n")

	b.WriteString("## ✅ Covered\n\n")
	hasCovered := false
	for _, g := range groups {
		if len(g.Tests) == 0 {
			continue
		}
		hasCovered = true
		writeGroup(&b, g)
	}
	if !hasCovered {
		b.WriteString("_No groups are currently covered._\n\n")
	}

	b.WriteString("## ❌ Not yet covered\n\n")
	hasMissing := false
	for _, g := range groups {
		if len(g.Tests) > 0 {
			continue
		}
		hasMissing = true
		writeGroup(&b, g)
	}
	if !hasMissing {
		b.WriteString("_All registered groups are covered. 🎉_\n\n")
	}

	return b.String()
}

func writeGroup(b *strings.Builder, g groupCoverage) {
	fmt.Fprintf(b, "### `%s`\n\n", g.TFName)
	if g.Category != "" {
		fmt.Fprintf(b, "_Category: %s_\n\n", g.Category)
	}
	// Only render the first non-empty line; ResourceGroup.Description often
	// embeds a redundant "Available operations" dump that is already covered
	// by the CRUD table below.
	if summary := firstLine(g.Description); summary != "" {
		fmt.Fprintf(b, "%s\n\n", summary)
	}
	if len(g.Ops) > 0 {
		b.WriteString("| CRUD | Method | Path |\n| --- | --- | --- |\n")
		for _, op := range g.Ops {
			fmt.Fprintf(b, "| %s | `%s` | `%s` |\n", op.CRUD, op.Method, op.Path)
		}
		b.WriteString("\n")
	}
	if len(g.Tests) > 0 {
		b.WriteString("Tests:")
		for _, t := range g.Tests {
			fmt.Fprintf(b, " `%s`", t)
		}
		b.WriteString("\n\n")
	}
}

func percent(n, total int) int {
	if total == 0 {
		return 0
	}
	return n * 100 / total
}

// firstLine returns the first non-empty trimmed line of s.
func firstLine(s string) string {
	for _, line := range strings.Split(s, "\n") {
		if t := strings.TrimSpace(line); t != "" {
			return t
		}
	}
	return ""
}
