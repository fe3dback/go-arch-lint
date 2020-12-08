package validator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
)

type (
	utils struct {
		pathResolver  PathResolver
		document      spec.Document
		rootDirectory string
	}
)

func newUtils(
	pathResolver PathResolver,
	spec spec.Document,
	rootDirectory string,
) *utils {
	return &utils{
		pathResolver:  pathResolver,
		document:      spec,
		rootDirectory: rootDirectory,
	}
}

func (u *utils) assertVendorImportPathValid(importPath string) error {
	localPath := fmt.Sprintf("vendor/%s", importPath)
	err := u.assertPathValid(localPath)
	if err != nil {
		return fmt.Errorf("vendor dep '%s' not installed, run 'go mod vendor' first: %v",
			importPath,
			err,
		)
	}

	return nil
}

func (u *utils) assertPathValid(localPath string) error {
	absPath := filepath.Clean(fmt.Sprintf("%s/%s", u.rootDirectory, localPath))
	resolved, err := u.pathResolver.Resolve(absPath)
	if err != nil {
		return fmt.Errorf("failed to resolv path: %v", err)
	}

	if len(resolved) == 0 {
		return fmt.Errorf("not found directories for '%s' in '%s'", localPath, absPath)
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
	for knownName := range u.document.Components {
		if name == knownName {
			return nil
		}
	}

	return fmt.Errorf("unknown component '%s'", name)
}

func (u *utils) assertKnownVendor(name string) error {
	for knownName := range u.document.Vendors {
		if name == knownName {
			return nil
		}
	}

	return fmt.Errorf("unknown vendor '%s'", name)
}