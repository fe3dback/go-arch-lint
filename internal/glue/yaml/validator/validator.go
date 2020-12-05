package validator

import (
	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type Validator struct {
	refResolver   YamlSourceCodeReferenceResolver
	pathResolver  PathResolver
	rootDirectory string
}

func NewValidator(
	refResolver YamlSourceCodeReferenceResolver,
	pathResolver PathResolver,
	rootDirectory string,
) *Validator {
	return &Validator{
		refResolver:   refResolver,
		pathResolver:  pathResolver,
		rootDirectory: rootDirectory,
	}
}

func (v *Validator) Validate(doc spec.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	utils := newUtils(v.pathResolver, doc, v.rootDirectory)
	validators := []validator{
		newValidatorCommonComponents(v.refResolver, utils),
		newValidatorCommonVendors(v.refResolver, utils),
		newValidatorComponents(v.refResolver, utils),
		newValidatorDeps(v.refResolver, utils),
		newValidatorDepsComponents(v.refResolver, utils),
		newValidatorDepsVendors(v.refResolver, utils),
		newValidatorExcludeFiles(v.refResolver),
		newValidatorVendors(v.refResolver, utils),
		newValidatorVersion(v.refResolver),
	}

	for _, specValidator := range validators {
		notices = append(notices, specValidator.Validate(doc)...)
	}

	return notices
}
