package resolver

import (
	"context"
	"fmt"
	"path"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
)

type Resolver struct {
	projectFilesResolver projectFilesResolver
	projectFilesHolder   projectFilesHolder
}

func NewResolver(
	projectFilesResolver projectFilesResolver,
	projectFilesHolder projectFilesHolder,
) *Resolver {
	return &Resolver{
		projectFilesResolver: projectFilesResolver,
		projectFilesHolder:   projectFilesHolder,
	}
}

func (r *Resolver) ProjectFiles(ctx context.Context, spec arch.Spec) ([]models.FileHold, error) {
	scanDirectory := path.Clean(fmt.Sprintf("%s/%s",
		spec.RootDirectory.Value,
		spec.WorkingDirectory.Value,
	))

	projectFiles, err := r.projectFilesResolver.Scan(
		ctx,
		scanDirectory,
		spec.ModuleName.Value,
		refPathToList(spec.Exclude),
		refRegExpToList(spec.ExcludeFilesMatcher),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve project files: %w", err)
	}

	holdFiles := r.projectFilesHolder.HoldProjectFiles(projectFiles, spec.Components)
	return holdFiles, nil
}
