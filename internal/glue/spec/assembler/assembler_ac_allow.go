package assembler

import (
	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
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

func (efa allowAssembler) assemble(spec *speca.Spec, yamlSpec *spec.Document) error {
	spec.Allow = speca.Allow{
		DepOnAnyVendor: speca.NewReferableBool(
			yamlSpec.Allow.DepOnAnyVendor,
			efa.provideYamlRef("$.allow.depOnAnyVendor"),
		),
	}

	return nil
}
