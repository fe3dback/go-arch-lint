package assembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	Assembler struct {
		decoder      archDecoder
		validator    archValidator
		pathResolver pathResolver
	}
)

func NewAssembler(
	decoder archDecoder,
	validator archValidator,
	pathResolver pathResolver,
) *Assembler {
	return &Assembler{
		decoder:      decoder,
		validator:    validator,
		pathResolver: pathResolver,
	}
}

func (sa *Assembler) Assemble(prj common.Project) (speca.Spec, error) {
	spec := speca.Spec{
		RootDirectory: common.NewEmptyReferable(prj.Directory),
		ModuleName:    common.NewEmptyReferable(prj.ModuleName),
		Integrity: speca.Integrity{
			DocumentNotices: []speca.Notice{},
			Suggestions:     []speca.Notice{},
		},
	}

	document, schemeNotices, err := sa.decoder.Decode(prj.GoArchFilePath)
	if err != nil {
		return spec, fmt.Errorf("failed to decode document '%s': %w", prj.GoArchFilePath, err)
	}

	if len(schemeNotices) > 0 {
		// only simple scheme validation errors
		spec.Integrity.DocumentNotices = append(spec.Integrity.DocumentNotices, schemeNotices...)
	} else {
		// if scheme is ok, need check arch errors
		advancedErrors := sa.validator.Validate(document)
		spec.Integrity.DocumentNotices = append(spec.Integrity.DocumentNotices, advancedErrors...)
	}

	if document == nil {
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

	err = assembler.assemble(&spec, document)
	if err != nil {
		return spec, fmt.Errorf("failed to assemble document: %w", err)
	}

	return spec, nil
}
