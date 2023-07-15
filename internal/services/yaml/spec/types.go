package spec

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	yamlSourceCodeReferenceResolver interface {
		Resolve(filePath string, yamlPath string) models.Reference
	}

	jsonSchemaProvider interface {
		Provide(version int) ([]byte, error)
	}
)
