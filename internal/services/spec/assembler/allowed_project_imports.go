package assembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
)

type allowedProjectImportsAssembler struct {
	projectPath string
	resolver    *resolver
}

func newAllowedProjectImportsAssembler(
	projectPath string,
	resolver *resolver,
) *allowedProjectImportsAssembler {
	return &allowedProjectImportsAssembler{
		projectPath: projectPath,
		resolver:    resolver,
	}
}

func (aia *allowedProjectImportsAssembler) assemble(
	yamlDocument arch.Document,
	componentNames []string,
) ([]models.ResolvedPath, error) {
	list := make([]models.ResolvedPath, 0)

	allowedComponents := make([]string, 0)
	allowedComponents = append(allowedComponents, componentNames...)
	for _, componentName := range yamlDocument.CommonComponents().List() {
		allowedComponents = append(allowedComponents, componentName.Value())
	}

	for _, name := range allowedComponents {
		yamlComponent, ok := yamlDocument.Components().Map()[name]
		if !ok {
			continue
		}

		for _, componentIn := range yamlComponent.RelativePaths() {
			relativeGlobPath := componentIn.Value()

			resolved, err := aia.resolver.resolveLocalGlobPath(string(relativeGlobPath))
			if err != nil {
				return nil, fmt.Errorf("failed to resolve component path '%s'", relativeGlobPath)
			}

			list = append(list, resolved...)
		}
	}

	return list, nil
}
