package check

import (
	"context"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	SpecAssembler interface {
		Assemble() (speca.Spec, error)
	}

	ReferenceRender interface {
		SourceCode(ref models.Reference, height int, highlight bool) []byte
	}

	SpecChecker interface {
		Check(ctx context.Context, spec speca.Spec) (models.CheckResult, error)
	}
)
