package handlers

import (
	"fmt"
	"os"
	"strings"
)

// UploadFile is one file of a multipart upload: the local file sent as the
// form field named after its path in the repository.
type UploadFile struct {
	RepoPath  string
	LocalPath string
}

// ParseUploadFiles parses --file values of the form <repository path>=<local file>.
// The value is split at the first '=', so a local path may contain '='. Each
// local file must exist and be a regular file.
func ParseUploadFiles(values []string) ([]UploadFile, error) {
	files := make([]UploadFile, 0, len(values))
	for _, v := range values {
		repoPath, localPath, ok := strings.Cut(v, "=")
		repoPath = strings.TrimPrefix(repoPath, "/")
		if !ok || repoPath == "" || localPath == "" {
			return nil, fmt.Errorf("invalid --file %q: expected <repository path>=<local file>", v)
		}
		info, err := os.Stat(localPath)
		if err != nil {
			return nil, fmt.Errorf("invalid --file %q: %w", v, err)
		}
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("invalid --file %q: %s is not a regular file", v, localPath)
		}
		files = append(files, UploadFile{RepoPath: repoPath, LocalPath: localPath})
	}
	return files, nil
}
