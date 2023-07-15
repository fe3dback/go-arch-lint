package resolver

import (
	"context"
	"regexp"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	projectFilesResolver interface {
		Scan(
			ctx context.Context,
			projectDirectory string,
			moduleName string,
			excludePaths []models.ResolvedPath,
			excludeFileMatchers []*regexp.Regexp,
		) ([]models.ProjectFile, error)
	}

	projectFilesHolder interface {
		HoldProjectFiles(files []models.ProjectFile, components []speca.Component) []models.FileHold
	}
)
