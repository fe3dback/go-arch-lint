package specassembler

import (
	"fmt"
	"regexp"

	yaml "github.com/fe3dback/go-arch-lint/internal/glue/yamlspecprovider"
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type excludeFilesMatcherAssembler struct {
}

func newExcludeFilesMatcherAssembler() *excludeFilesMatcherAssembler {
	return &excludeFilesMatcherAssembler{}
}

func (efa excludeFilesMatcherAssembler) assemble(spec *models.ArchSpec, yamlSpec *yaml.YamlSpec) error {
	for _, regString := range yamlSpec.ExcludeFilesRegExp {
		matcher, err := regexp.Compile(regString)
		if err != nil {
			return fmt.Errorf("failed to compile regular expression '%s': %v", regString, err)
		}

		spec.ExcludeFilesMatcher = append(spec.ExcludeFilesMatcher, matcher)
	}

	return nil
}
