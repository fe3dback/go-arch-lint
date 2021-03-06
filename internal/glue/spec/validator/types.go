package validator

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	validator interface {
		Validate(doc arch.Document) []speca.Notice
	}

	YamlSourceCodeReferenceResolver interface {
		Resolve(yamlPath string) models.Reference
	}

	PathResolver interface {
		Resolve(absPath string) (resolvePaths []string, err error)
	}
)
