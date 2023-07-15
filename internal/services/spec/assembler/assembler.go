package assembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	Assembler struct {
		provider     archProvider
		validator    archValidator
		pathResolver pathResolver
	}
)

func NewAssembler(
	provider archProvider,
	validator archValidator,
	pathResolver pathResolver,
) *Assembler {
	return &Assembler{
		provider:     provider,
		validator:    validator,
		pathResolver: pathResolver,
	}
}

func (sa *Assembler) Assemble(prj common.Project) (speca.Spec, error) {
	spec := speca.Spec{
		RootDirectory: speca.NewEmptyReferable(prj.Directory),
		ModuleName:    speca.NewEmptyReferable(prj.ModuleName),
		Integrity: speca.Integrity{
			DocumentNotices: []speca.Notice{},
			Suggestions:     []speca.Notice{},
		},
	}

	yamlSpec, schemeNotices, err := sa.provider.Provide(prj.GoArchFilePath)
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
		prj.Directory,
		prj.ModuleName,
	)

	assembler := newSpecCompositeAssembler([]assembler{
		newComponentsAssembler(
			resolver,
			newAllowedProjectImportsAssembler(
				prj.Directory,
				resolver,
			),
			newAllowedVendorImportsAssembler(
				resolver,
			),
		),
		newExcludeAssembler(resolver),
		newExcludeFilesMatcherAssembler(),
		newAllowAssembler(),
		newWorkdirAssembler(),
	})

	err = assembler.assemble(&spec, yamlSpec)
	if err != nil {
		return spec, fmt.Errorf("failed to assemble yamlSpec: %w", err)
	}

	return spec, nil
}
