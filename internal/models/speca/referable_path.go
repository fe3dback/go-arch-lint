package speca

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	ReferableResolvedPath struct {
		value models.ResolvedPath
		ref   models.Reference
	}
)

func NewReferableResolvedPath(value models.ResolvedPath, ref models.Reference) ReferableResolvedPath {
	return ReferableResolvedPath{value: value, ref: ref}
}

func NewEmptyReferableResolvedPath(value models.ResolvedPath) ReferableResolvedPath {
	return ReferableResolvedPath{value: value, ref: NewEmptyReference()}
}

func (s ReferableResolvedPath) Reference() models.Reference {
	return s.ref
}

func (s ReferableResolvedPath) Value() models.ResolvedPath {
	return s.value
}
