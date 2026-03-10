// Package output provides rendering helpers for CLI output.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
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
// and single structs by reflecting over exported fields.
func renderTable(w io.Writer, v any) error {
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
