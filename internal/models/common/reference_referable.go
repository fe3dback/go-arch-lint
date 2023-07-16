package common

type Referable[T any] struct {
	Value     T
	Reference Reference
}

func NewReferable[T any](value T, ref Reference) Referable[T] {
	return Referable[T]{Value: value, Reference: ref}
}

func NewEmptyReferable[T any](value T) Referable[T] {
	return Referable[T]{Value: value, Reference: NewEmptyReference()}
}
