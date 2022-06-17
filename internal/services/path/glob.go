package path

import (
	"os"
	"path/filepath"
	"strings"
)

// COPYRIGHT: https://github.com/yargevad
// SOURCE: https://github.com/yargevad/filepathx/blob/907099cb5a626c26c0b04e01adca2798e63d030e/filepathx.go

// globs represents one filepath glob, with its elements joined by "**".
type globs []string

// Glob adds double-star support to the core path/filepath Glob function.
// It's useful when your globs might have double-stars, but you're not sure.
func glob(pattern string) ([]string, error) {
	if !strings.Contains(pattern, "**") {
		// passthru to core package if no double-star
		return filepath.Glob(pattern)
	}
	return globs(strings.Split(pattern, "**")).expand()
}

// expand finds matches for the provided globs.
func (globs globs) expand() ([]string, error) {
	matches := []string{""} // accumulate here
	for _, glob := range globs {
		var hits []string
		hitMap := map[string]bool{}
		for _, match := range matches {
			paths, err := filepath.Glob(match + glob)
			if err != nil {
				return nil, err
			}
			for _, path := range paths {
				err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					// save deduped match from current iteration
					if _, ok := hitMap[path]; !ok {
						hits = append(hits, path)
						hitMap[path] = true
					}
					return nil
				})
				if err != nil {
					return nil, err
				}
			}
		}
		matches = hits
	}

	// fix up return value for nil input
	if globs == nil && len(matches) > 0 && matches[0] == "" {
		matches = matches[1:]
	}

	return matches, nil
}
