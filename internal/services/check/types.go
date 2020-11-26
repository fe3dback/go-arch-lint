package check

import (
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	SpecAssembler interface {
		Assemble() (speca.Spec, error)
	}
)
