package checker

import (
	"context"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	projectFilesResolver interface {
		ProjectFiles(ctx context.Context, spec speca.Spec) ([]models.FileHold, error)
	}

	checker interface {
		Check(ctx context.Context, spec speca.Spec) (models.CheckResult, error)
	}

	sourceCodeRenderer interface {
		SourceCodeWithoutOffset(ref models.CodeReference, highlight bool) []byte
	}
)
