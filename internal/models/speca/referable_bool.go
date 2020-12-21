package speca

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	ReferableBool struct {
		value bool
		ref   models.Reference
	}
)

func NewReferableBool(value bool, ref models.Reference) ReferableBool {
	return ReferableBool{value: value, ref: ref}
}

func NewEmptyReferableBool(value bool) ReferableBool {
	return ReferableBool{value: value, ref: NewEmptyReference()}
}

func (s ReferableBool) Reference() models.Reference {
	return s.ref
}

func (s ReferableBool) Value() bool {
	return s.value
}
