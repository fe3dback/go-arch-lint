package xpath

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type (
	typeMatcher interface {
		match(ctx *queryContext, query models.FileQuery) ([]models.FileDescriptor, error)
	}

	fileScanner interface {
		Scan(scanDirectory string, fn func(path string, isDir bool) error) error
	}
)
