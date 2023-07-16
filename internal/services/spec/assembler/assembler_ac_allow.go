package assembler

import (
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type allowAssembler struct {
}

func newAllowAssembler() *allowAssembler {
	return &allowAssembler{}
}

func (efa *allowAssembler) assemble(spec *speca.Spec, document spec.Document) error {
	spec.Allow = speca.Allow{
		DepOnAnyVendor: document.Options().IsDependOnAnyVendor(),
	}

	return nil
}
