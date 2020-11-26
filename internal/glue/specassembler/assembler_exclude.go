package specassembler

import (
	"fmt"

	yaml "github.com/fe3dback/go-arch-lint/internal/glue/yamlspecprovider"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type excludeAssembler struct {
	resolver       *resolver
	provideYamlRef provideYamlRef
}

func newExcludeAssembler(
	resolver *resolver,
	provideYamlRef provideYamlRef,
) *excludeAssembler {
	return &excludeAssembler{
		resolver:       resolver,
		provideYamlRef: provideYamlRef,
	}
}

func (ea excludeAssembler) assemble(spec *speca.Spec, yamlSpec *yaml.YamlSpec) error {
	for _, yamlRelativePath := range yamlSpec.Exclude {
		tmpResolvedPath, err := ea.resolver.resolveLocalPath(yamlRelativePath)
		if err != nil {
			return fmt.Errorf("failed to assemble exclude '%s' path's: %v", yamlRelativePath, err)
		}

		resolvedPath := wrapPaths(ea.provideYamlRef("$.exclude"), tmpResolvedPath)
		spec.Exclude = append(spec.Exclude, resolvedPath...)
	}

	return nil
}
