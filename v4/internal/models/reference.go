package models

type Ref[T any] struct {
	Value T
	Ref   Reference
}

func NewRef[T any](value T, ref Reference) Ref[T] {
	return Ref[T]{
		Value: value,
		Ref:   ref,
	}
}

type RefSlice[T any] []Ref[T]

type Reference struct {
	File   PathAbsolute `json:"File"`
	Line   int          `json:"Line"`
	Column int          `json:"Column"`
	Valid  bool         `json:"Valid"`
}

func NewReference(File PathAbsolute, Line int, Column int) Reference {
	return Reference{
		File:   File,
		Line:   Line,
		Column: Column,
		Valid:  true,
	}
}

func NewInvalidReference() Reference {
	return Reference{
		Valid: false,
	}
}

type RefMap[K comparable, V any] struct {
	Values map[K]V
	Refs   map[K]Reference
}

func NewRefMap[K comparable, V any](size int) RefMap[K, V] {
	return RefMap[K, V]{
		Values: make(map[K]V, size),
		Refs:   make(map[K]Reference, size),
	}
}

func RefMapFrom[K comparable, V any](values map[K]V, refs map[K]Reference) RefMap[K, V] {
	return RefMap[K, V]{
		Values: values,
		Refs:   refs,
	}
}

func (rf *RefMap[K, V]) Set(key K, val V, ref Reference) {
	rf.Values[key] = val
	rf.Refs[key] = ref
}
