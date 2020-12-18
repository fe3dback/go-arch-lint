package assembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	assembler interface {
		assemble(spec *speca.Spec, yamlSpec *spec.Document) error
	}

	specModifier struct {
		modifiers []assembler
	}
)

func newSpecAssembler(modifiers []assembler) *specModifier {
	return &specModifier{
		modifiers: modifiers,
	}
}

func (s specModifier) assemble(spec *speca.Spec, yamlSpec *spec.Document) error {
	for _, modifier := range s.modifiers {
		err := modifier.assemble(spec, yamlSpec)
		if err != nil {
			return fmt.Errorf("failed to assemble spec with '%T' assembler", modifier)
		}
	}

	return nil
}
