package specassembler

import (
	yaml "github.com/fe3dback/go-arch-lint/internal/glue/yamlspecprovider"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	YamlSpecProvider interface {
		Provide() (*yaml.YamlSpec, error)
	}

	YamlSourceCodeReferenceResolver interface {
		Resolve(yamlPath string) speca.Reference
	}

	PathResolver interface {
		Resolve(absPath string) (resolvePaths []string, err error)
	}
)
