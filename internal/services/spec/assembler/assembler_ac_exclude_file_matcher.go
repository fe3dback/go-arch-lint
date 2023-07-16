package assembler

import (
	"regexp"

	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type excludeFilesMatcherAssembler struct{}

func newExcludeFilesMatcherAssembler() *excludeFilesMatcherAssembler {
	return &excludeFilesMatcherAssembler{}
}

func (efa *excludeFilesMatcherAssembler) assemble(spec *speca.Spec, yamlSpec spec.Document) error {
	for _, regString := range yamlSpec.ExcludedFilesRegExp().List() {
		matcher, err := regexp.Compile(regString.Value)
		if err != nil {
			continue
		}

		spec.ExcludeFilesMatcher = append(spec.ExcludeFilesMatcher, common.NewReferable(
			matcher,
			regString.Reference,
		))
	}

	return nil
}
