package path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type (
	Resolver struct{}
)

func NewResolver() *Resolver {
	return &Resolver{}
}

func (r Resolver) Resolve(absPath string) (resolvePaths []string, err error) {
	absPath = strings.TrimSuffix(absPath, ".")

	matches, err := glob(absPath)
	if err != nil {
		return nil, fmt.Errorf("can`t match path mask '%s': %w", absPath, err)
	}

	dirs := make([]string, 0)
	for _, match := range matches {
		fileInfo, err := os.Stat(match)
		if err != nil {
			return nil, fmt.Errorf("nostat '%s': %w", match, err)
		}

		switch mode := fileInfo.Mode(); {
		case mode.IsDir():
			// glob can return path with mixed separators on windows,
			// normalize it to os specific form
			dirs = append(dirs, filepath.Clean(match))
		default:
			continue
		}
	}

	return dirs, nil
}
