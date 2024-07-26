package validator

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type internalValidator interface {
	Validate(ctx *validationContext)
}

type pathHelper interface {
	MatchProjectFiles(somePath any) ([]models.FileRef, error)
}
