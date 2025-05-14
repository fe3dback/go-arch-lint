package validator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type (
	utils struct {
		pathResolver pathResolver
		document     spec.Document
		projectDir   string
	}
)

func newUtils(
	pathResolver pathResolver,
	document spec.Document,
	projectDir string,
) *utils {
	return &utils{
		projectDir:   projectDir,
		pathResolver: pathResolver,
		document:     document,
	}
}

func (u *utils) assertGlobPathValid(localGlobPath string) error {
	absPath := filepath.Join(u.projectDir, localGlobPath)
	resolved, err := u.pathResolver.Resolve(absPath)
	if err != nil {
		return fmt.Errorf("failed to resolv path: %w", err)
	}

	if len(resolved) == 0 && !u.document.Options().IgnoreNotFoundComponents().Value {
		return fmt.Errorf("not found directories for '%s' in '%s'", localGlobPath, absPath)
	}

	return u.assertDirectoriesValid(resolved...)
}

func (u *utils) assertDirectoriesValid(paths ...string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("directory '%s' not exist", path)
		}
	}

	return nil
}

func (u *utils) assertKnownComponent(name string) error {
	for knownName := range u.document.Components() {
		if name == knownName {
			return nil
		}
	}

	return fmt.Errorf("unknown component '%s'", name)
}

func (u *utils) assertKnownVendor(name string) error {
	for knownName := range u.document.Vendors() {
		if name == knownName {
			return nil
		}
	}

	return fmt.Errorf("unknown vendor '%s'", name)
}
