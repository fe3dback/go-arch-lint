package codeprinter

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type (
	LinesExtractor interface {
		ExtractLines(file arch.PathAbsolute, from int, to int) ([]string, error)
	}
)
