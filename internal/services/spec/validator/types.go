package validator

import (
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type (
	validator interface {
		Validate(doc spec.Document) []arch.Notice
	}

	pathResolver interface {
		Resolve(absPath string) (resolvePaths []string, err error)
	}
)
