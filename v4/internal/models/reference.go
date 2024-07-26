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
	values map[K]V
	refs   map[K]Reference
}

func NewRefMap[K comparable, V any](size int) RefMap[K, V] {
	return RefMap[K, V]{
		values: make(map[K]V, size),
		refs:   make(map[K]Reference, size),
	}
}

func (rf *RefMap[K, V]) Len() int {
	return len(rf.values)
}

func (rf *RefMap[K, V]) Set(key K, val V, ref Reference) {
	rf.values[key] = val
	rf.refs[key] = ref
}

func (rf *RefMap[K, V]) Get(key K) (V, Reference, bool) {
	value, hasValue := rf.values[key]
	ref, hasRef := rf.refs[key]

	return value, ref, hasValue && hasRef
}

func (rf *RefMap[K, V]) Has(key K) bool {
	_, hasValue := rf.values[key]
	return hasValue
}

func (rf *RefMap[K, V]) Each(fn func(K, V, Reference)) {
	for k, v := range rf.values {
		ref, exist := rf.refs[k]
		if !exist {
			continue
		}

		fn(k, v, ref)
	}
}
