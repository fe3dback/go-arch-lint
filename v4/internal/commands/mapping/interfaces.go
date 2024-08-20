package mapping

import "github.com/fe3dback/go-arch-lint-sdk/arch"

type (
	specFetcher interface {
		Fetch() (arch.Spec, error)
	}
)
