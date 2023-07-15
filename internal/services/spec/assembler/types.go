package assembler

import (
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	archProvider interface {
		Provide(archFile string) (arch.Document, []speca.Notice, error)
	}

	archValidator interface {
		Validate(doc arch.Document) []speca.Notice
	}

	pathResolver interface {
		Resolve(absPath string) (resolvePaths []string, err error)
	}
)
