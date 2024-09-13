package xstdout

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/pkg/codeprinter"
)

type (
	codePrinter interface {
		Print(ref arch.Reference, opts codeprinter.CodePrintOpts) (string, error)
	}
)
