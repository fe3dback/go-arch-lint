package assembler

import (
	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	YamlSpecProvider interface {
		Provide() (spec.ArchDocument, error)
	}

	YamlSourceCodeReferenceResolver interface {
		Resolve(yamlPath string) models.Reference
	}

	PathResolver interface {
		Resolve(absPath string) (resolvePaths []string, err error)
	}

	Validator interface {
		Validate(spec speca.Spec) []speca.Notice
	}
)
