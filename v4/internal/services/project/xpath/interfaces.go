package xpath

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type (
	typeMatcher interface {
		match(ctx *queryContext, query models.FileQuery) ([]models.FileDescriptor, error)
	}
)
