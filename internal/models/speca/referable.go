package speca

import "github.com/fe3dback/go-arch-lint/internal/models"

type Referable[T any] struct {
	value T
	ref   models.Reference
}

func NewReferable[T any](value T, ref models.Reference) Referable[T] {
	return Referable[T]{value: value, ref: ref}
}

func NewEmptyReferable[T any](value T) Referable[T] {
	return Referable[T]{value: value, ref: NewEmptyReference()}
}

func (s Referable[T]) Reference() models.Reference {
	return s.ref
}

func (s Referable[T]) Value() T {
	return s.value
}
