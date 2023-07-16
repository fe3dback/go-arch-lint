package assembler

import (
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type workdirAssembler struct {
}

func newWorkdirAssembler() *workdirAssembler {
	return &workdirAssembler{}
}

func (efa *workdirAssembler) assemble(spec *speca.Spec, document spec.Document) error {
	spec.WorkingDirectory = document.WorkingDirectory()

	return nil
}
