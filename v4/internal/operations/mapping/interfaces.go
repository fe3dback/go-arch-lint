package mapping

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type specFetcher interface {
	FetchSpec() (models.Spec, error)
}
