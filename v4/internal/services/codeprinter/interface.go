package codeprinter

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type (
	LinesExtractor interface {
		ExtractLines(file models.PathAbsolute, from int, to int) ([]string, error)
	}
)
