package validator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint/spec/archfile"
)

type (
	validatorUtils struct {
		pathResolver  PathResolver
		rootDirectory string
		spec          *archfile.YamlSpec
	}
)

func newValidatorUtils(
	pathResolver PathResolver,
	spec *archfile.YamlSpec,
	rootDirectory string,
) *validatorUtils {
	return &validatorUtils{
		pathResolver:  pathResolver,
		rootDirectory: rootDirectory,
		spec:          spec,
	}
}

func (u *validatorUtils) isValidImportPath(importPath string) error {
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

func (u *validatorUtils) isValidPath(localPath string) error {
	absPath := filepath.Clean(fmt.Sprintf("%s/%s", u.rootDirectory, localPath))
	resolved, err := u.pathResolver.Resolve(absPath)
	if err != nil {
		return fmt.Errorf("failed to resolv path: %v", err)
	}

	if len(resolved) == 0 {
		return fmt.Errorf("not found directories for '%s' in '%s'", localPath, absPath)
	}

	return u.isValidDirectories(resolved...)
}

func (u *validatorUtils) isValidDirectories(paths ...string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("directory '%s' not exist", path)
		}
	}

	return nil
}

func (u *validatorUtils) isKnownComponent(name string) error {
	for knownName := range u.spec.Components {
		if name == knownName {
			return nil
		}
	}

	return fmt.Errorf("unknown component '%s'", name)
}

func (u *validatorUtils) isKnownVendor(name string) error {
	for knownName := range u.spec.Vendors {
		if name == knownName {
			return nil
		}
	}

	return fmt.Errorf("unknown vendor '%s'", name)
}
