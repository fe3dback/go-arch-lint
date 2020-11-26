package specassembler

import (
	yaml "github.com/fe3dback/go-arch-lint/internal/glue/yamlspecprovider"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type allowAssembler struct {
	provideYamlRef provideYamlRef
}

func newAllowAssembler(provideYamlRef provideYamlRef) *allowAssembler {
	return &allowAssembler{
		provideYamlRef: provideYamlRef,
	}
}

func (efa allowAssembler) assemble(spec *speca.Spec, yamlSpec *yaml.YamlSpec) error {
	spec.Allow = speca.Allow{
		DepOnAnyVendor: speca.NewReferableBool(
			yamlSpec.Allow.DepOnAnyVendor,
			efa.provideYamlRef("$.allow.depOnAnyVendor"),
		),
	}

	return nil
}
