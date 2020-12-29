package assembler

import (
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type workdirAssembler struct {
}

func newWorkdirAssembler() *workdirAssembler {
	return &workdirAssembler{}
}

func (efa workdirAssembler) assemble(spec *speca.Spec, document arch.Document) error {
	spec.WorkingDirectory = document.WorkingDirectory()

	return nil
}
