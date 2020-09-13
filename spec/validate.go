package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	pathresolv "github.com/fe3dback/go-arch-lint/path"

	"github.com/goccy/go-yaml"
)

type (
	validateFn func() error
	validator  struct {
		rootDirectory string
		spec          YamlSpec
		source        []byte
		warnings      []string
	}
)

func newValidator(spec YamlSpec, source []byte, rootDirectory string) *validator {
	return &validator{
		rootDirectory: rootDirectory,
		spec:          spec,
		source:        source,
		warnings:      []string{},
	}
}

func (v *validator) validate() error {
	// validators
	v.validateVersion()
	v.validateComponents()
	v.validateExcludeFiles()
	v.validateDeps()
	v.validateVendors()
	v.validateCommonComponents()
	v.validateCommonVendors()

	// display warnings
	if warnings, ok := v.getWarnings(); !ok {
		for _, warning := range warnings {
			fmt.Printf("[Archfile] %s\n", warning)
		}

		return fmt.Errorf("syntax error")
	}

	return nil
}

func (v *validator) validateVersion() {
	v.check("$.version", func() error {
		if v.spec.Version <= supportedVersion && v.spec.Version > 0 {
			return nil
		}

		return fmt.Errorf("version %d is not supported, supported: [%d]",
			v.spec.Version,
			supportedVersion,
		)
	})
}

func (v *validator) validateComponents() {
	v.check("$.components", func() error {
		if len(v.spec.Components) == 0 {
			return fmt.Errorf("at least one component should by defined")
		}

		for name, component := range v.spec.Components {
			v.check(fmt.Sprintf("$.components.%s.in", name), func() error {
				return v.isValidPath(component.LocalPath)
			})
		}

		return nil
	})
}

func (v *validator) validateExcludeFiles() {
	for index, regExp := range v.spec.ExcludeFilesRegExp {
		v.check(fmt.Sprintf("$.excludeFiles[%d]", index), func() error {
			_, err := regexp.Compile(regExp)
			return err
		})
	}
}

func (v *validator) validateDeps() {
	for name, rules := range v.spec.Dependencies {
		v.check(fmt.Sprintf("$.deps.%s", name), func() error {
			return v.isKnownComponent(name)
		})

		for index, componentName := range rules.MayDependOn {
			v.check(fmt.Sprintf("$.deps.%s.mayDependOn[%d]", name, index), func() error {
				return v.isKnownComponent(componentName)
			})
		}

		for index, vendorName := range rules.CanUse {
			v.check(fmt.Sprintf("$.deps.%s.canUse[%d]", name, index), func() error {
				return v.isKnownVendor(vendorName)
			})
		}

		if len(rules.MayDependOn) == 0 && len(rules.CanUse) == 0 {
			v.check(fmt.Sprintf("$.deps.%s", name), func() error {
				if rules.AnyProjectDeps {
					return nil
				}

				if rules.AnyVendorDeps {
					return nil
				}

				return fmt.Errorf("should have ref in 'MayDependOn' or at least one flag of ['anyProjectDeps', 'anyVendorDeps']")
			})
		}
	}
}

func (v *validator) validateVendors() {
	for name, vendor := range v.spec.Vendors {
		v.check(fmt.Sprintf("$.vendors.%s.in", name), func() error {
			return v.isValidImportPath(vendor.ImportPath)
		})
	}
}

func (v *validator) validateCommonComponents() {
	for index, componentName := range v.spec.CommonComponents {
		v.check(fmt.Sprintf("$.commonComponents[%d]", index), func() error {
			return v.isKnownComponent(componentName)
		})
	}
}

func (v *validator) validateCommonVendors() {
	for index, vendorName := range v.spec.CommonVendors {
		v.check(fmt.Sprintf("$.commonVendors[%d]", index), func() error {
			return v.isKnownVendor(vendorName)
		})
	}
}

func (v *validator) check(path string, fn validateFn) {
	defer func() {
		if err := recover(); err != nil {
			v.warnings = append(v.warnings, fmt.Sprintf("not found path '%s': %v", path, err))
			return
		}
	}()

	checkError := fn()
	if checkError == nil {
		return
	}

	sourceLine, err := yaml.PathString(path)
	if err != nil {
		v.warnings = append(v.warnings, fmt.Sprintf("failed check '%s': %v", path, err))
		return
	}

	highlightSource, err := sourceLine.AnnotateSource(v.source, true)
	if err != nil {
		v.warnings = append(v.warnings, fmt.Sprintf("failed annotate '%s': %v", path, err))
		return
	}

	warnings := fmt.Sprintf("path '%s': %v\n%s", path, checkError, highlightSource)
	v.warnings = append(v.warnings, warnings)
}

func (v *validator) getWarnings() ([]string, bool) {
	if len(v.warnings) == 0 {
		return nil, true
	}

	return v.warnings, false
}

func (v *validator) isValidImportPath(importPath string) error {
	localPath := fmt.Sprintf("vendor/%s", importPath)
	err := v.isValidPath(localPath)
	if err != nil {
		return fmt.Errorf("vendor dep '%s' not installed, run 'go mod vendor' first: %v",
			importPath,
			err,
		)
	}

	return nil
}

func (v *validator) isValidPath(localPath string) error {
	absPath := filepath.Clean(fmt.Sprintf("%s/%s", v.rootDirectory, localPath))
	resolved, err := pathresolv.ResolvePath(absPath)
	if err != nil {
		return fmt.Errorf("failed to resolv path: %v", err)
	}

	if len(resolved) == 0 {
		return fmt.Errorf("not found directories for '%s' in '%s'", localPath, absPath)
	}

	return v.isValidDirectories(resolved...)
}

func (v *validator) isValidDirectories(paths ...string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("directory '%s' not exist", path)
		}
	}

	return nil
}

func (v *validator) isKnownComponent(name string) error {
	for knownName := range v.spec.Components {
		if name == knownName {
			return nil
		}
	}

	return fmt.Errorf("unknown component '%s'", name)
}

func (v *validator) isKnownVendor(name string) error {
	for knownName := range v.spec.Vendors {
		if name == knownName {
			return nil
		}
	}

	return fmt.Errorf("unknown vendor '%s'", name)
}
