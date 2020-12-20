package assembler

import (
	"fmt"
	"regexp"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type excludeFilesMatcherAssembler struct {
	provideYamlRef provideYamlRef
}

func newExcludeFilesMatcherAssembler(provideYamlRef provideYamlRef) *excludeFilesMatcherAssembler {
	return &excludeFilesMatcherAssembler{
		provideYamlRef: provideYamlRef,
	}
}

func (efa excludeFilesMatcherAssembler) assemble(spec *speca.Spec, yamlSpec *spec.ArchV1Document) error {
	for ind, regString := range yamlSpec.V1ExcludeFilesRegExp {
		ref := efa.provideYamlRef(fmt.Sprintf("$.excludeFiles[%d]", ind))

		matcher, err := regexp.Compile(regString)
		if err != nil {
			continue
		}

		spec.ExcludeFilesMatcher = append(spec.ExcludeFilesMatcher, speca.NewReferableRegExp(
			matcher,
			ref,
		))
	}

	return nil
}
