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
	"time"

	"github.com/fatih/color"
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
	case "markdown":
		return renderMarkdown(w, v)
	case "id":
		return renderIDs(w, v)
	default:
		return fmt.Errorf("unknown output format: %s (valid: table, markdown, json, id)", Format)
	}
}

// padRight pads s with spaces to width n.
func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}

// mdTable writes a markdown-style table to w given headers and rows.
func mdTable(w io.Writer, headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}
	// Compute column widths (based on plain text length, not ANSI codes).
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(stripANSI(h))
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				plain := len(stripANSI(cell))
				if plain > widths[i] {
					widths[i] = plain
				}
			}
		}
	}
	// Header row.
	parts := make([]string, len(headers))
	for i, h := range headers {
		parts[i] = padWithANSI(h, widths[i])
	}
	fmt.Fprintf(w, "| %s |\n", strings.Join(parts, " | "))
	// Separator row.
	seps := make([]string, len(headers))
	for i := range headers {
		seps[i] = strings.Repeat("-", widths[i])
	}
	fmt.Fprintf(w, "| %s |\n", strings.Join(seps, " | "))
	// Data rows.
	for _, row := range rows {
		cells := make([]string, len(headers))
		for i := range headers {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			cells[i] = padWithANSI(val, widths[i])
		}
		fmt.Fprintf(w, "| %s |\n", strings.Join(cells, " | "))
	}
}

// stripANSI removes ANSI escape sequences for width calculation.
func stripANSI(s string) string {
	var out strings.Builder
	inEsc := false
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEsc = false
			}
			continue
		}
		out.WriteRune(r)
	}
	return out.String()
}

// padWithANSI pads a string that may contain ANSI codes to a visual width.
func padWithANSI(s string, width int) string {
	plainLen := len(stripANSI(s))
	if plainLen >= width {
		return s
	}
	return s + strings.Repeat(" ", width-plainLen)
}

// renderTable renders v as a human-friendly terminal table (key-value blocks).
func renderTable(w io.Writer, v any) error {
	if items, ok := v.([]any); ok {
		return renderMapSliceKV(w, items)
	}
	if m, ok := v.(map[string]any); ok {
		return renderMapKV(w, m)
	}

	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Slice:
		if val.Len() == 0 {
			fmt.Fprintln(w, "(no results)")
			return nil
		}
		elem := val.Index(0)
		for elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		headers, indices := tableRowsFor(elem.Type())
		for i := range val.Len() {
			if i > 0 {
				fmt.Fprintln(w)
			}
			item := val.Index(i)
			for item.Kind() == reflect.Ptr {
				item = item.Elem()
			}
			writeKVBlock(w, headers, func(j int) string {
				return fmt.Sprintf("%v", deref(item.Field(indices[j])))
			})
		}
	case reflect.Struct:
		t := val.Type()
		var keys []string
		for i := range t.NumField() {
			f := t.Field(i)
			if f.IsExported() {
				keys = append(keys, f.Name)
			}
		}
		writeKVBlock(w, keys, func(j int) string {
			return fmt.Sprintf("%v", deref(val.FieldByName(keys[j])))
		})
	default:
		fmt.Fprintln(w, fmt.Sprintf("%v", v))
	}
	return nil
}

// renderMapSliceKV renders a []any (slice of maps) as key-value blocks.
func renderMapSliceKV(w io.Writer, items []any) error {
	if len(items) == 0 {
		fmt.Fprintln(w, "(no results)")
		return nil
	}
	first, ok := items[0].(map[string]any)
	if !ok {
		for _, item := range items {
			fmt.Fprintf(w, "%v\n", item)
		}
		return nil
	}
	cols := pickMapColumns(first)
	for i, item := range items {
		if i > 0 {
			fmt.Fprintln(w)
		}
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		writeKVBlock(w, cols, func(j int) string {
			val := flatValue(m[cols[j]])
			if cols[j] == "state" {
				val = colorState(val)
			}
			return val
		})
	}
	return nil
}

// renderMapKV renders a single map as key-value pairs.
func renderMapKV(w io.Writer, m map[string]any) error {
	cols := pickMapColumns(m)
	writeKVBlock(w, cols, func(j int) string {
		val := flatValue(m[cols[j]])
		if cols[j] == "state" {
			val = colorState(val)
		}
		return val
	})
	return nil
}

// writeKVBlock writes a block of key-value pairs with aligned keys.
func writeKVBlock(w io.Writer, keys []string, valueFn func(int) string) {
	bold := color.New(color.Bold)
	// Find the longest key for alignment.
	maxLen := 0
	for _, k := range keys {
		label := strings.ToUpper(k)
		if len(label) > maxLen {
			maxLen = len(label)
		}
	}
	for i, k := range keys {
		label := strings.ToUpper(k)
		fmt.Fprintf(w, "%s  %s\n", bold.Sprint(padRight(label, maxLen)), valueFn(i))
	}
}

// renderMarkdown renders v as a markdown-style table (for machine/AI consumption).
func renderMarkdown(w io.Writer, v any) error {
	if items, ok := v.([]any); ok {
		return renderMapSliceTable(w, items)
	}
	if m, ok := v.(map[string]any); ok {
		return renderMapTable(w, m)
	}

	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Slice:
		if val.Len() == 0 {
			fmt.Fprintln(w, "(no results)")
			return nil
		}
		elem := val.Index(0)
		for elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		headers, indices := tableRowsFor(elem.Type())
		bold := color.New(color.Bold)
		colorHeaders := make([]string, len(headers))
		for i, h := range headers {
			colorHeaders[i] = bold.Sprint(h)
		}
		var rows [][]string
		for i := range val.Len() {
			item := val.Index(i)
			for item.Kind() == reflect.Ptr {
				item = item.Elem()
			}
			fields := make([]string, len(indices))
			for j, idx := range indices {
				fields[j] = fmt.Sprintf("%v", deref(item.Field(idx)))
			}
			rows = append(rows, fields)
		}
		mdTable(w, colorHeaders, rows)
	case reflect.Struct:
		t := val.Type()
		bold := color.New(color.Bold)
		var headers []string
		var vals []string
		for i := range t.NumField() {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			headers = append(headers, bold.Sprint(f.Name))
			vals = append(vals, fmt.Sprintf("%v", deref(val.Field(i))))
		}
		for i, h := range headers {
			fmt.Fprintf(w, "%s: %s\n", h, vals[i])
		}
	default:
		fmt.Fprintln(w, fmt.Sprintf("%v", v))
	}
	return nil
}

// tableRowsFor returns the column headers and field indices for a struct type.
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

	used := make(map[int]bool)
	for i := range t.NumField() {
		name := t.Field(i).Name
		if !t.Field(i).IsExported() {
			continue
		}
		if order, ok := prioritySet[name]; ok {
			found = append(found, fieldEntry{name, i, order})
			used[i] = true
		}
	}

	for i := 1; i < len(found); i++ {
		for j := i; j > 0 && found[j].order < found[j-1].order; j-- {
			found[j], found[j-1] = found[j-1], found[j]
		}
	}

	// Append remaining exported fields in declaration order.
	nextOrder := len(priority)
	for i := range t.NumField() {
		if !t.Field(i).IsExported() || used[i] {
			continue
		}
		found = append(found, fieldEntry{t.Field(i).Name, i, nextOrder})
		nextOrder++
	}

	for _, f := range found {
		headers = append(headers, strings.ToUpper(f.name))
		indices = append(indices, f.index)
	}
	return
}

// renderIDs prints only ID fields for scripting.
func renderIDs(w io.Writer, v any) error {
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

// mapPriorityKeys controls which keys appear in table output for maps.
var mapPriorityKeys = []string{"id", "title", "state", "name", "display_name", "author", "user", "content", "source", "destination", "description", "created_on", "updated_on"}

// renderMapSliceTable renders a []any (slice of maps) as a markdown table.
func renderMapSliceTable(w io.Writer, items []any) error {
	if len(items) == 0 {
		fmt.Fprintln(w, "(no results)")
		return nil
	}

	first, ok := items[0].(map[string]any)
	if !ok {
		for _, item := range items {
			fmt.Fprintf(w, "%v\n", item)
		}
		return nil
	}

	cols := pickMapColumns(first)
	bold := color.New(color.Bold)
	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = bold.Sprint(strings.ToUpper(c))
	}

	var rows [][]string
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		fields := make([]string, len(cols))
		for i, c := range cols {
			val := flatValue(m[c])
			if c == "state" {
				val = colorState(val)
			}
			fields[i] = val
		}
		rows = append(rows, fields)
	}
	mdTable(w, headers, rows)
	return nil
}

// renderMapTable renders a single map[string]any as key-value pairs.
func renderMapTable(w io.Writer, m map[string]any) error {
	bold := color.New(color.Bold)
	cols := pickMapColumns(m)
	for _, k := range cols {
		val := flatValue(m[k])
		if k == "state" {
			val = colorState(val)
		}
		fmt.Fprintf(w, "%s: %s\n", bold.Sprint(strings.ToUpper(k)), val)
	}
	return nil
}

// pickMapColumns returns map keys to display: priority keys first (in defined
// order), then all remaining keys sorted alphabetically.
func pickMapColumns(m map[string]any) []string {
	seen := make(map[string]bool, len(m))
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
func flatValue(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		if d := formatDate(val); d != "" {
			return d
		}
		return sanitize(val)
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case map[string]any:
		return extractMapSummary(val)
	case []any:
		if len(val) == 0 {
			return ""
		}
		parts := make([]string, 0, len(val))
		for _, item := range val {
			if m, ok := item.(map[string]any); ok {
				parts = append(parts, extractMapSummary(m))
			} else {
				parts = append(parts, fmt.Sprintf("%v", item))
			}
		}
		return strings.Join(parts, ", ")
	default:
		return fmt.Sprintf("%v", val)
	}
}

// extractMapSummary pulls the most meaningful single value from a nested map.
func extractMapSummary(m map[string]any) string {
	for _, key := range []string{"display_name", "name", "title", "login", "username", "raw"} {
		if s, ok := m[key].(string); ok {
			return sanitize(s)
		}
	}
	if branch, ok := m["branch"].(map[string]any); ok {
		if name, ok := branch["name"].(string); ok {
			return name
		}
	}
	if href, ok := m["href"].(string); ok {
		return sanitize(href)
	}
	// Collect href values from sub-maps (e.g. links: {"html": {"href": "..."}}).
	var hrefs []string
	for _, v := range m {
		if sub, ok := v.(map[string]any); ok {
			if href, ok := sub["href"].(string); ok {
				hrefs = append(hrefs, href)
			}
		}
	}
	if len(hrefs) > 0 {
		sort.Strings(hrefs)
		return strings.Join(hrefs, " ")
	}
	b, _ := json.Marshal(m)
	return string(b)
}

// sanitize replaces newlines with spaces for single-line display.
func sanitize(s string) string {
	return strings.ReplaceAll(s, "\n", " ")
}

// truncate shortens s to maxLen, appending "…" if truncated.
func truncate(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}

// formatDate tries to parse s as an ISO 8601 timestamp and returns a
// human-readable representation. Returns "" if s is not a date.
func formatDate(s string) string {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Local().Format("02 Jan 2006 15:04")
		}
	}
	return ""
}

// colorState returns the state string with color applied.
func colorState(s string) string {
	switch strings.ToUpper(s) {
	case "OPEN":
		return color.GreenString(s)
	case "MERGED":
		return color.BlueString(s)
	case "DECLINED", "SUPERSEDED":
		return color.RedString(s)
	default:
		return s
	}
}

// deref dereferences a pointer reflect.Value for display.
func deref(v reflect.Value) any {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		return deref(v.Elem())
	}
	return v.Interface()
}
