package resolver

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type Resolver struct {
	projectFilesResolver ProjectFilesResolver
	projectFilesHolder   ProjectFilesHolder
}

func NewResolver(
	projectFilesResolver ProjectFilesResolver,
	projectFilesHolder ProjectFilesHolder,
) *Resolver {
	return &Resolver{
		projectFilesResolver: projectFilesResolver,
		projectFilesHolder:   projectFilesHolder,
	}
}

func (r *Resolver) ProjectFiles(
	rootDirectory string,
	moduleName string,
	spec speca.Spec,
) ([]models.FileHold, error) {
	projectFiles, err := r.projectFilesResolver.Scan(
		rootDirectory,
		moduleName,
		refPathToList(spec.Exclude),
		refRegExpToList(spec.ExcludeFilesMatcher),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve project files: %w", err)
	}

	holdFiles := r.projectFilesHolder.HoldProjectFiles(projectFiles, spec.Components)
	return holdFiles, nil
}
