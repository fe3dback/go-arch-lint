package assembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	Assembler struct {
		provider      ArchProvider
		validator     ArchValidator
		pathResolver  PathResolver
		rootDirectory string
		moduleName    string
	}
)

func NewAssembler(
	provider ArchProvider,
	validator ArchValidator,
	pathResolver PathResolver,
	rootDirectory string,
	moduleName string,
) *Assembler {
	return &Assembler{
		provider:      provider,
		validator:     validator,
		pathResolver:  pathResolver,
		rootDirectory: rootDirectory,
		moduleName:    moduleName,
	}
}

func (sa *Assembler) Assemble() (speca.Spec, error) {
	spec := speca.Spec{
		RootDirectory: speca.NewEmptyReferableString(sa.rootDirectory),
		ModuleName:    speca.NewEmptyReferableString(sa.moduleName),
		Integrity: speca.Integrity{
			DocumentNotices: []speca.Notice{},
			Suggestions:     []speca.Notice{},
		},
	}

	yamlSpec, schemeNotices, err := sa.provider.Provide()
	if err != nil {
		return spec, fmt.Errorf("failed to provide yamlSpec: %w", err)
	}

	if len(schemeNotices) > 0 {
		// only simple scheme validation errors
		spec.Integrity.DocumentNotices = append(spec.Integrity.DocumentNotices, schemeNotices...)
	} else {
		// if scheme is ok, need check arch errors
		advancedErrors := sa.validator.Validate(yamlSpec)
		spec.Integrity.DocumentNotices = append(spec.Integrity.DocumentNotices, advancedErrors...)
	}

	if yamlSpec == nil {
		return spec, nil
	}

	resolver := newResolver(
		sa.pathResolver,
		sa.rootDirectory,
		sa.moduleName,
	)

	assembler := newSpecCompositeAssembler([]assembler{
		newComponentsAssembler(
			resolver,
			newAllowedImportsAssembler(
				sa.rootDirectory,
				resolver,
			),
		),
		newExcludeAssembler(resolver),
		newExcludeFilesMatcherAssembler(),
		newAllowAssembler(),
	})

	err = assembler.assemble(&spec, yamlSpec)
	if err != nil {
		return spec, fmt.Errorf("failed to assemble yamlSpec: %w", err)
	}

	return spec, nil
}
