package xpath

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type FileScanner struct {
}

func NewFileScanner() *FileScanner {
	return &FileScanner{}
}

func (s *FileScanner) Scan(scanDirectory string, fn func(path string, isDir bool) error) error {
	return filepath.Walk(scanDirectory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed walking %q: %w", path, err)
		}

		return fn(path, info.IsDir())
	})
}
