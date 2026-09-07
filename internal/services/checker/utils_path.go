package checker

import (
	"path/filepath"
)

// relativeToRoot converts absolute file path to project relative path,
// always in unix form ('/internal/services/foo.go').
//
// go-arch-lint output (and json report) should not depend on host os,
// on windows filepath is stored with backslash separator, so it need to be normalized.
func relativeToRoot(rootDirectory string, absFilePath string) string {
	relativePath, err := filepath.Rel(rootDirectory, absFilePath)
	if err != nil {
		// path is outside of project root (should never happen), show it as is
		return filepath.ToSlash(absFilePath)
	}

	return "/" + filepath.ToSlash(relativePath)
}
