package assembler

import (
	"regexp"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type excludeFilesMatcherAssembler struct{}

func newExcludeFilesMatcherAssembler() *excludeFilesMatcherAssembler {
	return &excludeFilesMatcherAssembler{}
}

func (efa excludeFilesMatcherAssembler) assemble(spec *speca.Spec, yamlSpec arch.Document) error {
	for _, regString := range yamlSpec.ExcludedFilesRegExp().List() {
		matcher, err := regexp.Compile(regString.Value())
		if err != nil {
			continue
		}

		spec.ExcludeFilesMatcher = append(spec.ExcludeFilesMatcher, speca.NewReferableRegExp(
			matcher,
			regString.Reference(),
		))
	}

	return nil
}
