package speca

import (
	"regexp"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	ReferableRegExp struct {
		value *regexp.Regexp
		ref   models.Reference
	}
)

func NewReferableRegExp(value *regexp.Regexp, ref models.Reference) ReferableRegExp {
	return ReferableRegExp{value: value, ref: ref}
}

func NewEmptyReferableRegExp(value *regexp.Regexp) ReferableRegExp {
	return ReferableRegExp{value: value, ref: NewEmptyReference()}
}

func (s ReferableRegExp) Reference() models.Reference {
	return s.ref
}

func (s ReferableRegExp) Value() *regexp.Regexp {
	return s.value
}
