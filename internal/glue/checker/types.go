package checker

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	ProjectFilesResolver interface {
		ProjectFiles(spec speca.Spec) ([]models.FileHold, error)
	}
)
