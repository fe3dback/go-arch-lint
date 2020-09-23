package validator

import (
	"fmt"
	"os"
	"path/filepath"

	pathresolv "github.com/fe3dback/go-arch-lint/path"
	"github.com/fe3dback/go-arch-lint/spec/archfile"
)

type (
	utils struct {
		rootDirectory string
		spec          archfile.YamlSpec
	}
)

func newUtils(spec archfile.YamlSpec, rootDirectory string) *utils {
	return &utils{
		rootDirectory: rootDirectory,
		spec:          spec,
	}
}

func (u *utils) isValidImportPath(importPath string) error {
	localPath := fmt.Sprintf("vendor/%s", importPath)
	err := u.isValidPath(localPath)
	if err != nil {
		return fmt.Errorf("vendor dep '%s' not installed, run 'go mod vendor' first: %v",
			importPath,
			err,
		)
	}

	return nil
}

func (u *utils) isValidPath(localPath string) error {
	absPath := filepath.Clean(fmt.Sprintf("%s/%s", u.rootDirectory, localPath))
	resolved, err := pathresolv.ResolvePath(absPath)
	if err != nil {
		return fmt.Errorf("failed to resolv path: %v", err)
	}

	if len(resolved) == 0 {
		return fmt.Errorf("not found directories for '%s' in '%s'", localPath, absPath)
	}

	return u.isValidDirectories(resolved...)
}

func (u *utils) isValidDirectories(paths ...string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("directory '%s' not exist", path)
		}
	}

	return nil
}

func (u *utils) isKnownComponent(name string) error {
	for knownName := range u.spec.Components {
		if name == knownName {
			return nil
		}
	}

	return fmt.Errorf("unknown component '%s'", name)
}

func (u *utils) isKnownVendor(name string) error {
	for knownName := range u.spec.Vendors {
		if name == knownName {
			return nil
		}
	}

	return fmt.Errorf("unknown vendor '%s'", name)
}
