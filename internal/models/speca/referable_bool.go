package speca

type (
	ReferableBool struct {
		value bool
		ref   Reference
	}
)

func NewReferableBool(value bool, ref Reference) ReferableBool {
	return ReferableBool{value: value, ref: ref}
}

func NewEmptyReferableBool(value bool) ReferableBool {
	return ReferableBool{value: value, ref: NewEmptyReference()}
}

func (s ReferableBool) Reference() Reference {
	return s.ref
}

func (s ReferableBool) Value() bool {
	return s.value
}
