package assembler

import (
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type allowAssembler struct {
}

func newAllowAssembler() *allowAssembler {
	return &allowAssembler{}
}

func (efa allowAssembler) assemble(spec *speca.Spec, document arch.Document) error {
	spec.Allow = speca.Allow{
		DepOnAnyVendor: document.Options().IsDependOnAnyVendor(),
	}

	return nil
}
