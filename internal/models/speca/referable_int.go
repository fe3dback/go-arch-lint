package speca

type (
	ReferableInt struct {
		value int
		ref   Reference
	}
)

func NewReferableInt(value int, ref Reference) ReferableInt {
	return ReferableInt{value: value, ref: ref}
}

func NewEmptyReferableInt(value int) ReferableInt {
	return ReferableInt{value: value, ref: NewEmptyReference()}
}

func (s ReferableInt) Reference() Reference {
	return s.ref
}

func (s ReferableInt) Value() int {
	return s.value
}
