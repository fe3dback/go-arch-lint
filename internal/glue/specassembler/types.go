package specassembler

import (
	yaml "github.com/fe3dback/go-arch-lint/internal/glue/yamlspecprovider"
)

type (
	YamlSpecProvider interface {
		Provide() (*yaml.YamlSpec, error)
	}

	PathResolver interface {
		Resolve(absPath string) (resolvePaths []string, err error)
	}
)
