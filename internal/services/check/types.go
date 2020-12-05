package check

import (
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
)
