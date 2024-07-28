package models

type (
	Ref[T any] struct {
		Value T
		Ref   Reference
	}

	Reference struct {
		File   PathAbsolute `json:"File"`
		Line   int          `json:"Line"`
		Column int          `json:"Column"`
		XPath  string       `json:"-"`
		Valid  bool         `json:"Valid"`
	}
)

func NewRef[T any](value T, ref Reference) Ref[T] {
	return Ref[T]{
		Value: value,
		Ref:   ref,
	}
}

func NewReference(file PathAbsolute, line int, column int, xpath string) Reference {
	return Reference{
		File:   file,
		Line:   line,
		Column: column,
		XPath:  xpath,
		Valid:  true,
	}
}

func NewInvalidReference() Reference {
	return Reference{
		Valid: false,
	}
}

type RefSlice[T comparable] []Ref[T]

func (rs RefSlice[T]) Values() []T {
	list := make([]T, 0, len(rs))

	for _, refValue := range rs {
		list = append(list, refValue.Value)
	}

	return list
}

func (rs RefSlice[T]) Contains(ref Ref[T]) bool {
	for _, refValue := range rs {
		if refValue.Value == ref.Value {
			return true
		}
	}

	return false
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
