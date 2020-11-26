package specvalidator

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	validator interface {
		Validate(spec models.ArchSpec)
	}
)
