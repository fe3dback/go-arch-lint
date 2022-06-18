package selfInspect

import (
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	specAssembler interface {
		Assemble() (speca.Spec, error)
	}
)
