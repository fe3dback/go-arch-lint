package specassembler

import (
	"fmt"
	"regexp"

	yaml "github.com/fe3dback/go-arch-lint/internal/glue/yamlspecprovider"
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

func (efa excludeFilesMatcherAssembler) assemble(spec *speca.Spec, yamlSpec *yaml.YamlSpec) error {
	for ind, regString := range yamlSpec.ExcludeFilesRegExp {
		matcher, err := regexp.Compile(regString)
		if err != nil {
			return fmt.Errorf("failed to compile regular expression '%s': %v", regString, err)
		}

		spec.ExcludeFilesMatcher = append(spec.ExcludeFilesMatcher, speca.NewReferableRegExp(
			matcher,
			efa.provideYamlRef(fmt.Sprintf("$.excludeFiles[%d]", ind)),
		))
	}

	return nil
}
