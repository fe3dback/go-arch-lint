package specassembler

import (
	"fmt"

	yaml "github.com/fe3dback/go-arch-lint/internal/glue/yamlspecprovider"
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type excludeAssembler struct {
	resolver *resolver
}

func newExcludeAssembler(resolver *resolver) *excludeAssembler {
	return &excludeAssembler{
		resolver: resolver,
	}
}

func (ea excludeAssembler) assemble(spec *models.ArchSpec, yamlSpec *yaml.YamlSpec) error {
	for _, yamlRelativePath := range yamlSpec.Exclude {
		resolvedPath, err := ea.resolver.resolveLocalPath(yamlRelativePath)
		if err != nil {
			return fmt.Errorf("failed to assemble exclude '%s' path's: %v", yamlRelativePath, err)
		}

		spec.Exclude = append(spec.Exclude, resolvedPath...)
	}

	return nil
}
