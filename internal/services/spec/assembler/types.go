package assembler

import (
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type (
	archDecoder interface {
		Decode(archFile string) (spec.Document, []speca.Notice, error)
	}

	archValidator interface {
		Validate(doc spec.Document) []speca.Notice
	}

	pathResolver interface {
		Resolve(absPath string) (resolvePaths []string, err error)
	}
)
