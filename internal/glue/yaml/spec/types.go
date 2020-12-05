package spec

import (
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	Validator interface {
		Validate(doc Document) []speca.Notice
	}
)
