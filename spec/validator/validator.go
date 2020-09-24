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

	Warning struct {
		Path    string
		Warning error
	}
)

type (
	ArchFileValidator struct {
		registry *archFileCheckerRegistry
	}
)

func NewArchFileValidator(spec archfile.YamlSpec, rootDirectory string) *ArchFileValidator {
	return &ArchFileValidator{
		registry: newArchFileCheckerRegistry(
			spec,
			newUtils(spec, rootDirectory),
		),
	}
}

func (v *ArchFileValidator) Validate() []Warning {
	warnings := make([]Warning, 0)

	for _, checker := range v.registry.createdCheckers {
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
				Path:    checker.path,
				Warning: fmt.Errorf("not found path '%s': %v", checker.path, err),
			}

			return
		}
	}()

	if err := checker.checker(); err != nil {
		return &Warning{
			Path:    checker.path,
			Warning: fmt.Errorf("path '%s': %v", checker.path, err),
		}
	}

	return nil
}
