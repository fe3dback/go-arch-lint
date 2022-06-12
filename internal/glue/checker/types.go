package checker

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	projectFilesResolver interface {
		ProjectFiles(spec speca.Spec) ([]models.FileHold, error)
	}

	checker interface {
		Check(spec speca.Spec) (models.CheckResult, error)
	}

	sourceCodeRenderer interface {
		SourceCodeWithoutOffset(ref models.Reference, height int, highlight bool) []byte
	}
)
