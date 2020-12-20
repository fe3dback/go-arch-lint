package spec

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	YamlSourceCodeReferenceResolver interface {
		Resolve(yamlPath string) models.Reference
	}

	Validator interface {
		Validate(doc arch.Document) []speca.Notice
	}
)
