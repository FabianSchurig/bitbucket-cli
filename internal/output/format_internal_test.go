package output

import (
	"bytes"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/fatih/color"
)

const renderToErrFmt = "RenderTo: %v"

func TestRender_WritesToStdout(t *testing.T) {
	Format = "json"

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		done <- buf.String()
	}()

	if err := Render(map[string]any{"id": 7}); err != nil {
		t.Fatalf("Render: %v", err)
	}
	_ = w.Close()

	got := <-done
	if !strings.Contains(got, `"id": 7`) {
		t.Fatalf("expected stdout JSON output, got %q", got)
	}
}

func TestRenderTableStructFallback(t *testing.T) {
	Format = "table"
	type sample struct {
		ID      int
		Name    string
		private string
	}
	var buf bytes.Buffer
	if err := RenderTo(&buf, sample{ID: 1, Name: "demo", private: "hidden"}); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}
	got := buf.String()
	if !strings.Contains(got, "ID") || !strings.Contains(got, "NAME") || strings.Contains(got, "private") {
		t.Fatalf("unexpected struct table output:\n%s", got)
	}
}

func TestRenderTableDefaultFallback(t *testing.T) {
	Format = "table"
	var buf bytes.Buffer
	if err := RenderTo(&buf, 123); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}
	if got := strings.TrimSpace(buf.String()); got != "123" {
		t.Fatalf("expected fallback numeric output, got %q", got)
	}
}

func TestRenderTableAnySliceFallback(t *testing.T) {
	Format = "table"
	var buf bytes.Buffer
	if err := RenderTo(&buf, []any{"a", 2}); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}
	got := buf.String()
	if !strings.Contains(got, "a") || !strings.Contains(got, "2") {
		t.Fatalf("expected direct slice item output, got %q", got)
	}
}

func TestRenderMarkdownEmptySliceFallback(t *testing.T) {
	Format = "markdown"
	var buf bytes.Buffer
	if err := RenderTo(&buf, []sampleRow{}); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}
	if strings.TrimSpace(buf.String()) != noResults {
		t.Fatalf("expected %q, got %q", noResults, buf.String())
	}
}

func TestRenderMarkdownStructSliceFallback(t *testing.T) {
	Format = "markdown"
	var buf bytes.Buffer
	rows := []sampleRow{{Id: 1, Title: "One", State: "OPEN"}, {Id: 2, Title: "Two", State: "MERGED"}}
	if err := RenderTo(&buf, rows); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}
	got := buf.String()
	if !strings.Contains(got, "TITLE") || !strings.Contains(got, "One") || !strings.Contains(got, "Two") {
		t.Fatalf("unexpected markdown slice output:\n%s", got)
	}
}

func TestRenderMarkdownStructFallback(t *testing.T) {
	Format = "markdown"
	var buf bytes.Buffer
	if err := RenderTo(&buf, sampleRow{Id: 3, Title: "Three"}); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}
	got := buf.String()
	if !strings.Contains(got, "Id: 3") || !strings.Contains(got, "Title: Three") {
		t.Fatalf("unexpected markdown struct output:\n%s", got)
	}
}

func TestRenderMarkdownDefaultFallback(t *testing.T) {
	Format = "markdown"
	var buf bytes.Buffer
	if err := RenderTo(&buf, true); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}
	if strings.TrimSpace(buf.String()) != "true" {
		t.Fatalf("expected fallback boolean output, got %q", buf.String())
	}
}

func TestRenderMarkdownAnySliceFallback(t *testing.T) {
	Format = "markdown"
	var buf bytes.Buffer
	if err := RenderTo(&buf, []any{"alpha", "beta"}); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}
	got := buf.String()
	if !strings.Contains(got, "alpha") || !strings.Contains(got, "beta") {
		t.Fatalf("unexpected markdown any-slice output:\n%s", got)
	}
}

func TestRenderIDsMapSliceFallback(t *testing.T) {
	Format = "id"
	var buf bytes.Buffer
	items := []any{map[string]any{"id": float64(11)}, "skip", map[string]any{"id": float64(12)}}
	if err := RenderTo(&buf, items); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}
	got := buf.String()
	if !strings.Contains(got, "11") || !strings.Contains(got, "12") {
		t.Fatalf("expected map IDs, got %q", got)
	}
}

func TestRenderIDsSingleMapWithoutID(t *testing.T) {
	Format = "id"
	var buf bytes.Buffer
	if err := RenderTo(&buf, map[string]any{"title": "missing"}); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected empty output for map without id, got %q", buf.String())
	}
}

func TestRenderIDsSingleStructWithoutID(t *testing.T) {
	Format = "id"
	var buf bytes.Buffer
	if err := RenderTo(&buf, struct{ Title string }{Title: "x"}); err != nil {
		t.Fatalf(renderToErrFmt, err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected empty output for struct without Id, got %q", buf.String())
	}
}

func TestHelpersAndValueFormatting(t *testing.T) {
	color.NoColor = true
	t.Cleanup(func() { color.NoColor = false })

	if got := stripANSI("\033[32mOPEN\033[0m"); got != "OPEN" {
		t.Fatalf("expected ANSI stripped string, got %q", got)
	}

	if got := colorState("DECLINED"); got != "DECLINED" {
		t.Fatalf("expected declined state text, got %q", got)
	}
	if got := colorState("superseeded"); got != "superseeded" {
		t.Fatalf("expected superseded state text, got %q", got)
	}
	if got := colorIfState("state", "OPEN"); got != "OPEN" {
		t.Fatalf("expected colorIfState to preserve state text, got %q", got)
	}
	if got := colorIfState("title", "demo"); got != "demo" {
		t.Fatalf("expected non-state value unchanged, got %q", got)
	}

	values := map[string]string{
		"nil":         flatValue(nil),
		"float":       flatValue(12.5),
		"bool":        flatValue(true),
		"emptyArray":  flatValue([]any{}),
		"mixedArray":  flatValue([]any{"a", 2}),
		"defaultType": flatValue(struct{ Name string }{Name: "demo"}),
	}
	if values["nil"] != "" || values["float"] != "12.5" || values["bool"] != "true" || values["emptyArray"] != "" {
		t.Fatalf("unexpected flat values: %#v", values)
	}
	if values["mixedArray"] != "a, 2" {
		t.Fatalf("expected mixed array formatting, got %q", values["mixedArray"])
	}
	if !strings.Contains(values["defaultType"], "demo") {
		t.Fatalf("expected default type formatting to include value, got %q", values["defaultType"])
	}

	if got := extractMapSummary(map[string]any{"branch": map[string]any{"name": "feature/a"}}); got != "feature/a" {
		t.Fatalf("expected branch summary, got %q", got)
	}
	if got := extractMapSummary(map[string]any{"href": "https://example.test"}); got != "https://example.test" {
		t.Fatalf("expected href summary, got %q", got)
	}
	if got := extractMapSummary(map[string]any{
		"html": map[string]any{"href": "https://b.example"},
		"self": map[string]any{"href": "https://a.example"},
	}); got != "https://a.example https://b.example" {
		t.Fatalf("expected sorted nested hrefs, got %q", got)
	}
	if got := extractMapSummary(map[string]any{"count": float64(2)}); !strings.Contains(got, `"count":2`) {
		t.Fatalf("expected JSON fallback summary, got %q", got)
	}
}

func TestMarkdownHelpersAndErrorPaths(t *testing.T) {
	if err := mdTable(ioDiscardWithError{}, []string{"ID"}, [][]string{{"1"}}); !errors.Is(err, errWriteFailed) {
		t.Fatalf("expected markdown write error, got %v", err)
	}

	if err := mdTable(&bytes.Buffer{}, nil, nil); err != nil {
		t.Fatalf("expected empty headers to no-op, got %v", err)
	}
}

func TestWritePlainItemsAndPrintIDHelpers(t *testing.T) {
	if err := writePlainItems(ioDiscardWithError{}, []any{"x"}); !errors.Is(err, errWriteFailed) {
		t.Fatalf("expected plain item write error, got %v", err)
	}

	if err := printMapID(ioDiscardWithError{}, map[string]any{"id": 1}); !errors.Is(err, errWriteFailed) {
		t.Fatalf("expected printMapID write error, got %v", err)
	}

	if err := printFieldID(ioDiscardWithError{}, reflect.ValueOf(struct{ Id int }{Id: 1})); !errors.Is(err, errWriteFailed) {
		t.Fatalf("expected printFieldID write error, got %v", err)
	}
}

func TestRenderMapHelpersSkipNonMaps(t *testing.T) {
	var buf bytes.Buffer
	if err := renderMapSliceKV(&buf, []any{map[string]any{"id": 1}, "skip"}); err != nil {
		t.Fatalf("unexpected renderMapSliceKV error: %v", err)
	}

	buf.Reset()
	if err := renderMapSliceTable(&buf, []any{map[string]any{"id": 1}, "skip"}); err != nil {
		t.Fatalf("unexpected renderMapSliceTable error: %v", err)
	}
}

type sampleRow struct {
	Id    int
	Title string
	State string
}

var errWriteFailed = errors.New("write failed")

type ioDiscardWithError struct{}

func (ioDiscardWithError) Write(_ []byte) (int, error) {
	return 0, errWriteFailed
}
