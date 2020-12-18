package validator

import "github.com/fe3dback/go-arch-lint/internal/models/speca"

type (
	validator interface {
		Validate(spec speca.Spec) []speca.Notice
	}
)
