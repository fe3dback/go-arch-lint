package specassembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	provideYamlRef func(path string) speca.Reference

	SpecAssembler struct {
		provider              YamlSpecProvider
		pathResolver          PathResolver
		yamlReferenceResolver YamlSourceCodeReferenceResolver
		rootDirectory         string
		moduleName            string
	}
)

func NewSpecAssembler(
	provider YamlSpecProvider,
	pathResolver PathResolver,
	yamlReferenceResolver YamlSourceCodeReferenceResolver,
	rootDirectory string,
	moduleName string,
) *SpecAssembler {
	return &SpecAssembler{
		provider:              provider,
		pathResolver:          pathResolver,
		yamlReferenceResolver: yamlReferenceResolver,
		rootDirectory:         rootDirectory,
		moduleName:            moduleName,
	}
}

func (sa *SpecAssembler) Assemble() (speca.Spec, error) {
	spec := speca.Spec{
		RootDirectory: speca.NewEmptyReferableString(sa.rootDirectory),
		ModuleName:    speca.NewEmptyReferableString(sa.moduleName),
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

	yamlRefProvider := func(yamlPath string) speca.Reference {
		return sa.yamlReferenceResolver.Resolve(yamlPath)
	}

	assembler := newSpecAssembler([]assembler{
		newComponentsAssembler(
			resolver,
			newAllowedImportsAssembler(
				sa.rootDirectory,
				resolver,
			),
			yamlRefProvider,
		),
		newExcludeAssembler(resolver, yamlRefProvider),
		newExcludeFilesMatcherAssembler(yamlRefProvider),
		newAllowAssembler(yamlRefProvider),
	})

	err = assembler.assemble(&spec, yamlSpec)
	if err != nil {
		return spec, fmt.Errorf("failed to assemble yamlSpec: %w", err)
	}

	return spec, nil
}
