package assembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	assembler interface {
		assemble(spec *speca.Spec, doc arch.Document) error
	}

	specCompositeModifier struct {
		modifiers []assembler
	}
)

func newSpecCompositeAssembler(modifiers []assembler) *specCompositeModifier {
	return &specCompositeModifier{
		modifiers: modifiers,
	}
}

func (s *specCompositeModifier) assemble(spec *speca.Spec, doc arch.Document) error {
	for _, modifier := range s.modifiers {
		err := modifier.assemble(spec, doc)
		if err != nil {
			return fmt.Errorf("failed to assemble spec with '%T' assembler: %w", modifier, err)
		}
	}

	return nil
}
