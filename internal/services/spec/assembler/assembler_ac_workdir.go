package assembler

import (
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type workdirAssembler struct {
}

func newWorkdirAssembler() *workdirAssembler {
	return &workdirAssembler{}
}

func (efa *workdirAssembler) assemble(spec *arch.Spec, document spec.Document) error {
	spec.WorkingDirectory = document.WorkingDirectory()

	return nil
}
