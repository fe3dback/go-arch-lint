package assembler

import (
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	ArchProvider interface {
		Provide() (arch.Document, []speca.Notice, error)
	}

	ArchValidator interface {
		Validate(doc arch.Document) []speca.Notice
	}

	PathResolver interface {
		Resolve(absPath string) (resolvePaths []string, err error)
	}
)
