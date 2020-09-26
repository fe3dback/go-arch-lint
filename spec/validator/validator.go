package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/spec/archfile"
)

type (
	ArchFileValidatorFn func() error

	ArchFileRuleChecker struct {
		path    string
		checker ArchFileValidatorFn
	}
)

type (
	ArchFileValidator struct {
		rootDirectory string
		yamlSpec      *archfile.YamlSpec
	}
)

func NewArchFileValidator(yamlSpec *archfile.YamlSpec, rootDirectory string) *ArchFileValidator {
	return &ArchFileValidator{
		rootDirectory: rootDirectory,
		yamlSpec:      yamlSpec,
	}
}

func (v *ArchFileValidator) Validate() []Warning {
	utils := NewUtils(v.yamlSpec, v.rootDirectory)
	registry := newArchFileCheckerRegistry(v.yamlSpec, utils)
	warnings := make([]Warning, 0)

	for _, checker := range registry.createdCheckers {
		if warning := v.check(checker); warning != nil {
			warnings = append(warnings, *warning)
		}
	}

	return warnings
}

func (v *ArchFileValidator) check(checker ArchFileRuleChecker) (warn *Warning) {
	defer func() {
		if err := recover(); err != nil {
			warn = &Warning{
				yamlPath:    checker.path,
				yamlWarning: fmt.Errorf("not found path '%s': %v", checker.path, err),
			}

			return
		}
	}()

	if err := checker.checker(); err != nil {
		return &Warning{
			yamlPath:    checker.path,
			yamlWarning: fmt.Errorf("path '%s': %v", checker.path, err),
		}
	}

	return nil
}
