package validator

import (
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type Validator struct {
	pathResolver  PathResolver
	rootDirectory string
}

func NewValidator(
	pathResolver PathResolver,
	rootDirectory string,
) *Validator {
	return &Validator{
		pathResolver:  pathResolver,
		rootDirectory: rootDirectory,
	}
}

func (v *Validator) Validate(doc arch.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	utils := newUtils(v.pathResolver, doc, v.rootDirectory)
	validators := []validator{
		newValidatorCommonComponents(utils),
		newValidatorCommonVendors(utils),
		newValidatorComponents(utils),
		newValidatorDeps(utils),
		newValidatorDepsComponents(utils),
		newValidatorDepsVendors(utils),
		newValidatorExcludeFiles(),
		newValidatorVendors(utils),
		newValidatorVersion(),
	}

	for _, specValidator := range validators {
		notices = append(notices, specValidator.Validate(doc)...)
	}

	return notices
}
