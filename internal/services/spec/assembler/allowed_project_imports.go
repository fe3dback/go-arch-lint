package assembler

import (
	"fmt"
	"path"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
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
	yamlDocument spec.Document,
	componentNames []string,
) ([]models.ResolvedPath, error) {
	list := make([]models.ResolvedPath, 0)

	allowedComponents := make([]string, 0)
	allowedComponents = append(allowedComponents, componentNames...)
	for _, componentName := range yamlDocument.CommonComponents() {
		allowedComponents = append(allowedComponents, componentName.Value)
	}

	for _, name := range allowedComponents {
		yamlComponent, ok := yamlDocument.Components()[name]
		if !ok {
			continue
		}

		for _, componentIn := range yamlComponent.Value.RelativePaths() {
			relativeGlobPath := componentIn

			resolved, err := aia.resolver.resolveLocalGlobPath(
				path.Clean(fmt.Sprintf("%s/%s",
					yamlDocument.WorkingDirectory().Value,
					string(relativeGlobPath),
				)),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve component path '%s'", relativeGlobPath)
			}

			list = append(list, resolved...)
		}
	}

	return list, nil
}
