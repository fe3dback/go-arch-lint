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
	File   PathAbsolute
	Line   int
	Column int
	Valid  bool
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
	values map[K]V
	refs   map[K]Reference
}

func (rf *RefMap[K, V]) Value(key K) (V, bool) {
	val, exist := rf.values[key]
	return val, exist
}

func (rf *RefMap[K, V]) Reference(key K) (Reference, bool) {
	reference, exist := rf.refs[key]
	return reference, exist
}

func (rf *RefMap[K, V]) Data(key K) (V, Reference, bool) {
	val, existValue := rf.values[key]
	reference, existRef := rf.refs[key]

	return val, reference, existValue && existRef
}

func (rf *RefMap[K, V]) Insert(key K, val V, ref Reference) {
	rf.values[key] = val
	rf.refs[key] = ref
}
