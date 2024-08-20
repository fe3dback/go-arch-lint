package xstdout

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type (
	codePrinter interface {
		Print(ref arch.Reference, opts models.CodePrintOpts) (string, error)
	}
)
