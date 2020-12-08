package checker

import (
	"regexp"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	ProjectFilesResolver interface {
		Resolve(
			projectDirectory string,
			moduleName string,
			excludePaths []models.ResolvedPath,
			excludeFileMatchers []*regexp.Regexp,
		) ([]models.ResolvedFile, error)
	}
)
