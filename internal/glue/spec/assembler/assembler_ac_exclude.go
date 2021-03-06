package assembler

import (
	"fmt"
	"path"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type excludeAssembler struct {
	resolver *resolver
}

func newExcludeAssembler(
	resolver *resolver,
) *excludeAssembler {
	return &excludeAssembler{
		resolver: resolver,
	}
}

func (ea excludeAssembler) assemble(spec *speca.Spec, document arch.Document) error {
	for _, yamlRelativePath := range document.ExcludedDirectories().List() {
		tmpResolvedPath, err := ea.resolver.resolveLocalPath(
			path.Clean(fmt.Sprintf("%s/%s",
				document.WorkingDirectory().Value(),
				yamlRelativePath.Value(),
			)),
		)
		if err != nil {
			return fmt.Errorf("failed to assemble exclude '%s' path's: %w", yamlRelativePath.Value(), err)
		}

		resolvedPath := wrapPaths(yamlRelativePath.Reference(), tmpResolvedPath)
		spec.Exclude = append(spec.Exclude, resolvedPath...)
	}

	return nil
}
