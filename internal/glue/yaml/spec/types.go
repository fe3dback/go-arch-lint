package spec

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	YamlSourceCodeReferenceResolver interface {
		Resolve(yamlPath string) models.Reference
	}

	JsonSchemaProvider interface {
		Provide(version int) (string, error)
	}
)
