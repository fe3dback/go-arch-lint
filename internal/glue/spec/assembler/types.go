package assembler

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
)

type (
	ArchProvider interface {
		Provide() (arch.Arch, error)
	}

	YamlSourceCodeReferenceResolver interface {
		Resolve(yamlPath string) models.Reference
	}

	PathResolver interface {
		Resolve(absPath string) (resolvePaths []string, err error)
	}
)
