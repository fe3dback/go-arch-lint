package codeprinter

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type area struct {
	ref  arch.Reference
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
