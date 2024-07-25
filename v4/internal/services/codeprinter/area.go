package codeprinter

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type area struct {
	ref  models.Reference
	from int // inclusive
	to   int // inclusive
}

func (a area) lineNumbers(fn func(ind int, line int, isReferenced bool)) {
	ind := 0
	for line := a.from; line <= a.to; line++ {
		fn(ind, line, a.ref.Line == line)
		ind++
	}
}
