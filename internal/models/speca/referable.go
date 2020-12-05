package speca

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

func NewReference(file string, line int, offset int) models.Reference {
	return models.Reference{Valid: true, File: file, Line: line, Offset: offset}
}

func NewEmptyReference() models.Reference {
	return models.Reference{Valid: false, File: "", Line: 0, Offset: 0}
}
