package handlers

import (
	"bytes"
	"context"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
)

func writeTempFile(t *testing.T, dir, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
	return path
}

func TestParseUploadFiles(t *testing.T) {
	dir := t.TempDir()
	local := writeTempFile(t, dir, "a=b.docx", []byte("x"))

	got, err := ParseUploadFiles([]string{
		"suite x/tst_1/testdata/cover_page.docx=" + local,
		"/docs/readme.md=" + local,
	})
	if err != nil {
		t.Fatalf("ParseUploadFiles: %v", err)
	}
	want := []UploadFile{
		{RepoPath: "suite x/tst_1/testdata/cover_page.docx", LocalPath: local},
		{RepoPath: "docs/readme.md", LocalPath: local},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseUploadFiles = %#v, want %#v", got, want)
	}

	if got, err := ParseUploadFiles(nil); err != nil || len(got) != 0 {
		t.Errorf("ParseUploadFiles(nil) = %#v, %v; want empty, nil", got, err)
	}
}

func TestParseUploadFilesRejectsInvalidValues(t *testing.T) {
	dir := t.TempDir()
	local := writeTempFile(t, dir, "f.txt", []byte("x"))
	tests := []struct{ name, value, wantErr string }{
		{"missing separator", "docs/readme.md", "expected <repository path>=<local file>"},
		{"empty repository path", "=" + local, "expected <repository path>=<local file>"},
		{"empty local path", "docs/readme.md=", "expected <repository path>=<local file>"},
		{"missing local file", "docs/readme.md=" + filepath.Join(dir, "missing"), "no such file"},
		{"directory", "docs=" + dir, "is not a regular file"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseUploadFiles([]string{tt.value})
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ParseUploadFiles(%q) error = %v, want it to contain %q", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestDispatchRawSendsMultipartUpload(t *testing.T) {
	dir := t.TempDir()
	binary := []byte{0x50, 0x4b, 0x03, 0x04, 0x00, 0xff, '\r', '\n'}
	docx := writeTempFile(t, dir, "local.docx", binary)

	var (
		rawQuery string
		fields   = map[string][]string{}
		files    = map[string][]byte{}
		names    = map[string]string{}
		auth     string
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawQuery = r.URL.RawQuery
		auth = r.Header.Get("Authorization")
		mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mediaType != "multipart/form-data" {
			t.Errorf("Content-Type = %q, want multipart/form-data", r.Header.Get("Content-Type"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		reader := multipart.NewReader(r.Body, params["boundary"])
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Errorf("reading part: %v", err)
				break
			}
			content, _ := io.ReadAll(part)
			if part.FileName() != "" {
				files[part.FormName()] = content
				names[part.FormName()] = part.FileName()
			} else {
				fields[part.FormName()] = append(fields[part.FormName()], string(content))
			}
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	c := &client.BBClient{Client: newInternalTestClient(srv.URL).Client, Token: "secret"}
	_, err := DispatchRaw(context.Background(), c, Request{
		Method:      http.MethodPost,
		URLTemplate: srv.URL + "/repositories/{workspace}/{repo_slug}/src",
		PathParams:  map[string]string{"workspace": "ws", "repo_slug": "repo"},
		QueryParams: map[string]string{
			"message": "Add cover page",
			"branch":  "chore/cover-pages",
			"files":   "old/cover_page.docx",
			"author":  "",
		},
		Multipart: true,
		Files:     []UploadFile{{RepoPath: "suite x/tst_1/testdata/cover_page.docx", LocalPath: docx}},
	})
	if err != nil {
		t.Fatalf("DispatchRaw: %v", err)
	}

	if rawQuery != "" {
		t.Errorf("query string = %q, want the parameters in the form instead", rawQuery)
	}
	if auth != "Bearer secret" {
		t.Errorf("Authorization = %q, want the public API credentials", auth)
	}
	wantFields := map[string][]string{
		"message": {"Add cover page"},
		"branch":  {"chore/cover-pages"},
		"files":   {"old/cover_page.docx"},
	}
	if !reflect.DeepEqual(fields, wantFields) {
		t.Errorf("form fields = %#v, want %#v", fields, wantFields)
	}
	got, ok := files["suite x/tst_1/testdata/cover_page.docx"]
	if !ok || !bytes.Equal(got, binary) {
		t.Errorf("uploaded files = %#v, want the local file's bytes under its repository path", files)
	}
	if name := names["suite x/tst_1/testdata/cover_page.docx"]; name != "local.docx" {
		t.Errorf("file part name = %q, want the local file name", name)
	}
}
