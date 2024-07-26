package xstdout

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type (
	codePrinter interface {
		Print(ref models.Reference, opts models.CodePrintOpts) (string, error)
	}
)
