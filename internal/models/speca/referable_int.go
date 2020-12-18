package speca

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	ReferableInt struct {
		value int
		ref   models.Reference
	}
)

func NewReferableInt(value int, ref models.Reference) ReferableInt {
	return ReferableInt{value: value, ref: ref}
}

func NewEmptyReferableInt(value int) ReferableInt {
	return ReferableInt{value: value, ref: NewEmptyReference()}
}

func (s ReferableInt) Reference() models.Reference {
	return s.ref
}

func (s ReferableInt) Value() int {
	return s.value
}
