package validator

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type internalValidator interface {
	Validate(conf models.Config) error
}
