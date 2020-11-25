package specassembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type SpecAssembler struct {
	provider      YamlSpecProvider
	pathResolver  PathResolver
	rootDirectory string
	moduleName    string
}

func NewSpecAssembler(
	provider YamlSpecProvider,
	pathResolver PathResolver,
	rootDirectory string,
	moduleName string,
) *SpecAssembler {
	return &SpecAssembler{
		provider:      provider,
		pathResolver:  pathResolver,
		rootDirectory: rootDirectory,
		moduleName:    moduleName,
	}
}

func (sa *SpecAssembler) Assemble() (models.ArchSpec, error) {
	spec := models.ArchSpec{
		RootDirectory: sa.rootDirectory,
		ModuleName:    sa.moduleName,
	}

	yamlSpec, err := sa.provider.Provide()
	if err != nil {
		return spec, fmt.Errorf("failed to provide yamlSpec: %w", err)
	}

	resolver := newResolver(
		sa.pathResolver,
		sa.rootDirectory,
		sa.moduleName,
	)

	assembler := newSpecAssembler([]assembler{
		newComponentsAssembler(
			resolver,
			newAllowedImportsAssembler(
				sa.rootDirectory,
				resolver,
			),
		),
		newExcludeAssembler(resolver),
		newExcludeFilesMatcherAssembler(),
	})

	err = assembler.assemble(&spec, yamlSpec)
	if err != nil {
		return spec, fmt.Errorf("failed to assemble yamlSpec: %w", err)
	}

	return spec, nil
}
