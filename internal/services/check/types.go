package check

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	SpecAssembler interface {
		Assemble() (models.ArchSpec, error)
	}
)
