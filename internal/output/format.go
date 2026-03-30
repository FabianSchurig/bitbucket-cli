// Package output provides rendering helpers for CLI output.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"
)

// Format controls the output format. Overridden by the --output flag.
var Format = "table"

// Render writes v to stdout in the configured format.
func Render(v any) error {
	return RenderTo(os.Stdout, v)
}

// RenderTo writes v to the provided writer in the configured format.
func RenderTo(w io.Writer, v any) error {
	switch Format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(v)
	case "table":
		return renderTable(w, v)
	case "id":
		return renderIDs(w, v)
	default:
		return fmt.Errorf("unknown output format: %s (valid: json, table, id)", Format)
	}
}

// renderTable renders v as a tab-aligned table. It handles slices of structs
// and single structs by reflecting over exported fields, plus map[string]any
// and []any from generic API dispatch.
func renderTable(w io.Writer, v any) error {
	// Handle []any (e.g. collected paginated values from generic dispatch)
	if items, ok := v.([]any); ok {
		return renderMapSliceTable(w, items)
	}
	// Handle map[string]any (single resource from generic dispatch)
	if m, ok := v.(map[string]any); ok {
		return renderMapTable(w, m)
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Slice:
		if val.Len() == 0 {
			fmt.Fprintln(tw, "(no results)")
			return tw.Flush()
		}
		elem := val.Index(0)
		for elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		headers, _ := tableRowsFor(elem.Type())
		fmt.Fprintln(tw, strings.Join(headers, "\t"))
		for i := range val.Len() {
			item := val.Index(i)
			for item.Kind() == reflect.Ptr {
				item = item.Elem()
			}
			_, row := tableRowsFor(item.Type())
			fields := make([]string, len(row))
			for j, idx := range row {
				fields[j] = fmt.Sprintf("%v", deref(item.Field(idx)))
			}
			fmt.Fprintln(tw, strings.Join(fields, "\t"))
		}
	case reflect.Struct:
		t := val.Type()
		for i := range t.NumField() {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			fmt.Fprintf(tw, "%s\t%v\n", f.Name, deref(val.Field(i)))
		}
	default:
		fmt.Fprintln(tw, fmt.Sprintf("%v", v))
	}

	return tw.Flush()
}

// tableRowsFor returns the column headers and field indices for a struct type,
// limited to a predefined set of important fields for readability.
func tableRowsFor(t reflect.Type) (headers []string, indices []int) {
	priority := []string{"Id", "Title", "State", "Author", "Source", "Destination", "CreatedOn", "UpdatedOn"}
	prioritySet := make(map[string]int, len(priority))
	for i, p := range priority {
		prioritySet[p] = i
	}

	type fieldEntry struct {
		name  string
		index int
		order int
	}
	var found []fieldEntry

	for i := range t.NumField() {
		name := t.Field(i).Name
		if !t.Field(i).IsExported() {
			continue
		}
		if order, ok := prioritySet[name]; ok {
			found = append(found, fieldEntry{name, i, order})
		}
	}

	// Sort by priority order
	for i := 1; i < len(found); i++ {
		for j := i; j > 0 && found[j].order < found[j-1].order; j-- {
			found[j], found[j-1] = found[j-1], found[j]
		}
	}

	// If no priority fields found, use all exported fields
	if len(found) == 0 {
		for i := range t.NumField() {
			if t.Field(i).IsExported() {
				found = append(found, fieldEntry{t.Field(i).Name, i, i})
			}
		}
	}

	for _, f := range found {
		headers = append(headers, strings.ToUpper(f.name))
		indices = append(indices, f.index)
	}
	return
}

// renderIDs prints only ID fields (or the first field) for scripting.
func renderIDs(w io.Writer, v any) error {
	// Handle []any from generic dispatch
	if items, ok := v.([]any); ok {
		for _, item := range items {
			if m, ok := item.(map[string]any); ok {
				if id, ok := m["id"]; ok {
					fmt.Fprintln(w, flatValue(id))
				}
			}
		}
		return nil
	}
	// Handle map[string]any from generic dispatch
	if m, ok := v.(map[string]any); ok {
		if id, ok := m["id"]; ok {
			fmt.Fprintln(w, flatValue(id))
		}
		return nil
	}

	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		// Single object — print ID if available
		if f := val.FieldByName("Id"); f.IsValid() {
			fmt.Fprintln(w, deref(f))
		}
		return nil
	}

	for i := range val.Len() {
		item := val.Index(i)
		for item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		if f := item.FieldByName("Id"); f.IsValid() {
			fmt.Fprintln(w, deref(f))
		}
	}
	return nil
}

// mapPriorityKeys controls which keys appear first in table output for maps.
var mapPriorityKeys = []string{"id", "title", "state", "display_name", "name", "author", "created_on", "updated_on"}

// renderMapSliceTable renders a []any (slice of maps) as a table.
func renderMapSliceTable(w io.Writer, items []any) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	if len(items) == 0 {
		fmt.Fprintln(tw, "(no results)")
		return tw.Flush()
	}

	first, ok := items[0].(map[string]any)
	if !ok {
		// Not maps — fall back to one-value-per-line.
		for _, item := range items {
			fmt.Fprintf(tw, "%v\n", item)
		}
		return tw.Flush()
	}

	cols := pickMapColumns(first)
	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = strings.ToUpper(c)
	}
	fmt.Fprintln(tw, strings.Join(headers, "\t"))

	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		fields := make([]string, len(cols))
		for i, c := range cols {
			fields[i] = flatValue(m[c])
		}
		fmt.Fprintln(tw, strings.Join(fields, "\t"))
	}
	return tw.Flush()
}

// renderMapTable renders a single map[string]any as key-value pairs.
func renderMapTable(w io.Writer, m map[string]any) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	cols := pickMapColumns(m)
	for _, k := range cols {
		fmt.Fprintf(tw, "%s\t%s\n", strings.ToUpper(k), flatValue(m[k]))
	}
	return tw.Flush()
}

// pickMapColumns returns map keys in priority order, followed by remaining keys sorted.
func pickMapColumns(m map[string]any) []string {
	seen := make(map[string]bool)
	var cols []string
	for _, k := range mapPriorityKeys {
		if _, ok := m[k]; ok {
			cols = append(cols, k)
			seen[k] = true
		}
	}
	var rest []string
	for k := range m {
		if !seen[k] {
			rest = append(rest, k)
		}
	}
	sort.Strings(rest)
	cols = append(cols, rest...)
	return cols
}

// flatValue converts any value to a flat string for table display.
// Nested maps/slices are shown as compact JSON.
func flatValue(v any) string {
	if v == nil {
		return "<nil>"
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case map[string]any:
		// For nested objects, try to extract a meaningful single value.
		if name, ok := val["name"].(string); ok {
			return name
		}
		if dn, ok := val["display_name"].(string); ok {
			return dn
		}
		b, _ := json.Marshal(val)
		return string(b)
	case []any:
		b, _ := json.Marshal(val)
		return string(b)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// deref dereferences a pointer reflect.Value for display.
func deref(v reflect.Value) any {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "<nil>"
		}
		return deref(v.Elem())
	}
	return v.Interface()
}
