package resolver

import (
	"regexp"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	ProjectFilesResolver interface {
		Scan(
			projectDirectory string,
			moduleName string,
			excludePaths []models.ResolvedPath,
			excludeFileMatchers []*regexp.Regexp,
		) ([]models.ProjectFile, error)
	}

	ProjectFilesHolder interface {
		HoldProjectFiles(files []models.ProjectFile, components []speca.Component) []models.FileHold
	}
)
