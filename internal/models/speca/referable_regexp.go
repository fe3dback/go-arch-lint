package speca

import (
	"regexp"
)

type (
	ReferableRegExp struct {
		value *regexp.Regexp
		ref   Reference
	}
)

func NewReferableRegExp(value *regexp.Regexp, ref Reference) ReferableRegExp {
	return ReferableRegExp{value: value, ref: ref}
}

func NewEmptyReferableRegExp(value *regexp.Regexp) ReferableRegExp {
	return ReferableRegExp{value: value, ref: NewEmptyReference()}
}

func (s ReferableRegExp) Reference() Reference {
	return s.ref
}

func (s ReferableRegExp) Value() *regexp.Regexp {
	return s.value
}
