package xstdout

import (
	"github.com/fe3dback/go-arch-lint-sdk/pkg/codeprinter"
)

type (
	codePrinter interface {
		Print(ref codeprinter.Reference, opts codeprinter.CodePrintOpts) (string, error)
	}
)
