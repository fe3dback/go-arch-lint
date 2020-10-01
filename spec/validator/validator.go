package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/models"

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
		rootDirectory  string
		yamlSpec       *archfile.YamlSpec
		validatorUtils *validatorUtils
	}
)

func NewArchFileValidator(
	pathResolver PathResolver,
	yamlSpec *archfile.YamlSpec,
	rootDirectory string,
) *ArchFileValidator {
	return &ArchFileValidator{
		rootDirectory: rootDirectory,
		yamlSpec:      yamlSpec,
		validatorUtils: newValidatorUtils(
			pathResolver,
			yamlSpec,
			rootDirectory,
		),
	}
}

func (v *ArchFileValidator) Validate() []models.ArchFileSyntaxWarning {
	registry := newArchFileCheckerRegistry(v.yamlSpec, v.validatorUtils)
	warnings := make([]models.ArchFileSyntaxWarning, 0)

	for _, checker := range registry.createdCheckers {
		if warning := v.check(checker); warning != nil {
			warnings = append(warnings, *warning)
		}
	}

	return warnings
}

func (v *ArchFileValidator) check(checker ArchFileRuleChecker) (warn *models.ArchFileSyntaxWarning) {
	defer func() {
		if err := recover(); err != nil {
			warn = models.NewArchFileSyntaxWarning(
				checker.path,
				fmt.Errorf("not found path '%s': %v", checker.path, err),
			)

			return
		}
	}()

	if err := checker.checker(); err != nil {
		return models.NewArchFileSyntaxWarning(
			checker.path,
			fmt.Errorf("path '%s': %v", checker.path, err),
		)
	}

	return nil
}
