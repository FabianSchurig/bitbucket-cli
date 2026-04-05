package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/FabianSchurig/bitbucket-cli/internal/generated"
	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

const (
	renderToErrFmt = "RenderTo: %v"
	fixedDate      = "2026-01-01T00:00:00Z"
)

func TestRenderJSON(t *testing.T) {
	output.Format = "json"

	id := 42
	title := "Test PR"
	state := generated.PullrequestStateOPEN
	pr := generated.Pullrequest{
		Id:    &id,
		Title: &title,
		State: &state,
	}

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, pr); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}

	got := buf.String()
	if !strings.Contains(got, `"id": 42`) {
		t.Errorf("expected JSON to contain id=42, got: %s", got)
	}
	if !strings.Contains(got, `"title": "Test PR"`) {
		t.Errorf("expected JSON to contain title, got: %s", got)
	}
}

func TestRenderTable_EmptySlice(t *testing.T) {
	output.Format = "table"

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, []generated.Pullrequest{}); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}

	got := buf.String()
	if !strings.Contains(got, "(no results)") {
		t.Errorf("expected '(no results)', got: %s", got)
	}
}

func TestRenderTable_Slice(t *testing.T) {
	output.Format = "table"

	id1, id2 := 1, 2
	t1, t2 := "PR One", "PR Two"
	s1 := generated.PullrequestStateOPEN
	s2 := generated.PullrequestStateMERGED

	prs := []generated.Pullrequest{
		{Id: &id1, Title: &t1, State: &s1},
		{Id: &id2, Title: &t2, State: &s2},
	}

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, prs); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}

	got := buf.String()
	if !strings.Contains(got, "PR One") {
		t.Errorf("expected 'PR One' in table output, got: %s", got)
	}
	if !strings.Contains(got, "PR Two") {
		t.Errorf("expected 'PR Two' in table output, got: %s", got)
	}
}

func TestRenderID(t *testing.T) {
	output.Format = "id"

	id1, id2 := 10, 20
	t1, t2 := "A", "B"
	prs := []generated.Pullrequest{
		{Id: &id1, Title: &t1},
		{Id: &id2, Title: &t2},
	}

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, prs); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}

	got := buf.String()
	if !strings.Contains(got, "10") || !strings.Contains(got, "20") {
		t.Errorf("expected IDs 10 and 20, got: %s", got)
	}
}

func TestRenderUnknownFormat(t *testing.T) {
	output.Format = "xml"

	var buf bytes.Buffer
	err := output.RenderTo(&buf, "anything")
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestRenderMapSliceTable_AllKeysShown(t *testing.T) {
	output.Format = "table"

	items := []any{
		map[string]any{
			"id":         float64(1),
			"content":    map[string]any{"raw": "hello world"},
			"user":       map[string]any{"display_name": "alice"},
			"created_on": fixedDate,
			"updated_on": fixedDate,
			"pending":    false,
			"type":       "pullrequest_comment",
			"inline":     map[string]any{"path": "main.go"},
		},
	}

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, items); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}

	got := buf.String()

	// All keys must appear as column headers.
	for _, col := range []string{"ID", "CONTENT", "USER", "CREATED_ON", "UPDATED_ON", "INLINE", "PENDING", "TYPE"} {
		if !strings.Contains(got, col) {
			t.Errorf("expected column %q in table output, got:\n%s", col, got)
		}
	}
}

func TestRenderMapSliceTable_PriorityKeysFirst(t *testing.T) {
	output.Format = "table"

	items := []any{
		map[string]any{
			"id":         float64(1),
			"content":    map[string]any{"raw": "text"},
			"created_on": fixedDate,
			"zebra":      "extra",
			"alpha":      "extra",
		},
	}

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, items); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}

	got := buf.String()

	// In KV format each key is on its own line. Priority keys should come first.
	idPos := strings.Index(got, "ID")
	createdPos := strings.Index(got, "CREATED_ON")
	alphaPos := strings.Index(got, "ALPHA")
	zebraPos := strings.Index(got, "ZEBRA")

	if idPos < 0 || createdPos < 0 || alphaPos < 0 || zebraPos < 0 {
		t.Fatalf("missing expected keys in output:\n%s", got)
	}
	if idPos > createdPos {
		t.Errorf("ID should come before CREATED_ON")
	}
	if createdPos > alphaPos {
		t.Errorf("CREATED_ON (priority) should come before ALPHA (non-priority)")
	}
	if alphaPos > zebraPos {
		t.Errorf("ALPHA should come before ZEBRA (alphabetical)")
	}
}

func TestRenderMapTable_AllKeysShown(t *testing.T) {
	output.Format = "table"

	m := map[string]any{
		"id":     float64(42),
		"title":  "My PR",
		"custom": "value",
	}

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, m); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}

	got := buf.String()
	for _, col := range []string{"ID", "TITLE", "CUSTOM"} {
		if !strings.Contains(got, col) {
			t.Errorf("expected key %q in output, got:\n%s", col, got)
		}
	}
}

func TestFlatValue_FullTextNotTruncated(t *testing.T) {
	output.Format = "table"

	longText := strings.Repeat("a", 200)
	items := []any{
		map[string]any{
			"id":      float64(1),
			"content": map[string]any{"raw": longText},
		},
	}

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, items); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}

	got := buf.String()
	if !strings.Contains(got, longText) {
		t.Errorf("expected full text (200 chars) in output, but it was truncated:\n%s", got)
	}
}

func TestFlatValue_DateFormatting(t *testing.T) {
	output.Format = "table"

	items := []any{
		map[string]any{
			"id":         float64(1),
			"created_on": "2026-03-30T11:17:47.606820+00:00",
			"updated_on": "2026-01-15T09:30:00Z",
		},
	}

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, items); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}

	got := buf.String()
	if !strings.Contains(got, "30 Mar 2026") {
		t.Errorf("expected pretty-printed date '30 Mar 2026' in output, got:\n%s", got)
	}
	if !strings.Contains(got, "15 Jan 2026") {
		t.Errorf("expected pretty-printed date '15 Jan 2026' in output, got:\n%s", got)
	}
	// Should NOT contain raw ISO 8601 strings.
	if strings.Contains(got, "2026-03-30T") {
		t.Errorf("expected ISO 8601 date to be formatted, but raw date found:\n%s", got)
	}
}

func TestRenderTable_KVFormat(t *testing.T) {
	output.Format = "table"

	items := []any{
		map[string]any{"id": float64(1), "title": "First"},
		map[string]any{"id": float64(2), "title": "Second"},
	}

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, items); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}

	got := buf.String()
	// KV format should NOT contain markdown pipe separators.
	if strings.Contains(got, "| ") {
		t.Errorf("table format should not contain markdown pipes, got:\n%s", got)
	}
	// Should contain both records separated by a blank line.
	if !strings.Contains(got, "First") || !strings.Contains(got, "Second") {
		t.Errorf("expected both records in output, got:\n%s", got)
	}
	if !strings.Contains(got, "\n\n") {
		t.Errorf("expected blank line between records, got:\n%s", got)
	}
}

func TestRenderMarkdown_MapSlice(t *testing.T) {
	output.Format = "markdown"

	items := []any{
		map[string]any{"id": float64(1), "title": "First"},
		map[string]any{"id": float64(2), "title": "Second"},
	}

	var buf bytes.Buffer
	if err := output.RenderTo(&buf, items); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}

	got := buf.String()
	// Markdown format should contain pipe-delimited table.
	if !strings.Contains(got, "|") {
		t.Errorf("markdown format should contain pipes, got:\n%s", got)
	}
	if !strings.Contains(got, "---") {
		t.Errorf("markdown format should contain separator row, got:\n%s", got)
	}
	if !strings.Contains(got, "First") || !strings.Contains(got, "Second") {
		t.Errorf("expected both records in output, got:\n%s", got)
	}
}
