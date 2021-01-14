package spec

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	YAMLSourceCodeReferenceResolver interface {
		Resolve(yamlPath string) models.Reference
	}

	JSONSchemaProvider interface {
		Provide(version int) (string, error)
	}
)
