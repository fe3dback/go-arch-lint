package assembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	provideYamlRef func(path string) models.Reference

	Assembler struct {
		provider              ArchProvider
		pathResolver          PathResolver
		yamlReferenceResolver YamlSourceCodeReferenceResolver
		rootDirectory         string
		moduleName            string
	}
)

func NewAssembler(
	provider ArchProvider,
	pathResolver PathResolver,
	yamlReferenceResolver YamlSourceCodeReferenceResolver,
	rootDirectory string,
	moduleName string,
) *Assembler {
	return &Assembler{
		provider:              provider,
		pathResolver:          pathResolver,
		yamlReferenceResolver: yamlReferenceResolver,
		rootDirectory:         rootDirectory,
		moduleName:            moduleName,
	}
}

func (sa *Assembler) Assemble() (speca.Spec, error) {
	spec := speca.Spec{
		RootDirectory: speca.NewEmptyReferableString(sa.rootDirectory),
		ModuleName:    speca.NewEmptyReferableString(sa.moduleName),
		Integrity: speca.Integrity{
			DocumentNotices: []speca.Notice{},
			SpecNotices:     []speca.Notice{},
			Suggestions:     []speca.Notice{},
		},
	}

	archFile, err := sa.provider.Provide()
	if err != nil {
		return spec, fmt.Errorf("failed to provide yamlSpec: %w", err)
	}

	yamlSpec := archFile.Document()
	yamlIntegrity := archFile.Integrity()

	resolver := newResolver(
		sa.pathResolver,
		sa.rootDirectory,
		sa.moduleName,
	)

	yamlRefProvider := func(yamlPath string) models.Reference {
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

	err = assembler.assemble(&spec, &yamlSpec)
	if err != nil {
		return spec, fmt.Errorf("failed to assemble yamlSpec: %w", err)
	}

	// add integrity check in this level
	spec.Integrity.DocumentNotices = yamlIntegrity

	return spec, nil
}
