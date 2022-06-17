package mapping

import (
	"context"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	SpecAssembler interface {
		Assemble() (speca.Spec, error)
	}

	ProjectFilesResolver interface {
		ProjectFiles(ctx context.Context, spec speca.Spec) ([]models.FileHold, error)
	}
)
