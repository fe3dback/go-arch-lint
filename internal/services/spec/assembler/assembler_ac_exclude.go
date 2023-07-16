package assembler

import (
	"fmt"
	"path"

	"github.com/fe3dback/go-arch-lint/internal/models/speca"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
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

func (ea *excludeAssembler) assemble(spec *speca.Spec, document spec.Document) error {
	for _, yamlRelativePath := range document.ExcludedDirectories().List() {
		tmpResolvedPath, err := ea.resolver.resolveLocalGlobPath(
			path.Clean(fmt.Sprintf("%s/%s",
				document.WorkingDirectory().Value,
				yamlRelativePath.Value,
			)),
		)
		if err != nil {
			return fmt.Errorf("failed to assemble exclude '%s' path's: %w", yamlRelativePath.Value, err)
		}

		resolvedPath := wrap(yamlRelativePath.Reference, tmpResolvedPath)
		spec.Exclude = append(spec.Exclude, resolvedPath...)
	}

	return nil
}
