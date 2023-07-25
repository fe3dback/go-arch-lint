package decoder

import "github.com/fe3dback/go-arch-lint/internal/services/spec"

type doc interface {
	spec.Document

	postSetup()
}
