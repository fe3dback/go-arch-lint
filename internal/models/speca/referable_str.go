package speca

type (
	ReferableString struct {
		value string
		ref   Reference
	}
)

func NewReferableString(value string, ref Reference) ReferableString {
	return ReferableString{value: value, ref: ref}
}

func NewEmptyReferableString(value string) ReferableString {
	return ReferableString{value: value, ref: NewEmptyReference()}
}

func (s ReferableString) Reference() Reference {
	return s.ref
}

func (s ReferableString) Value() string {
	return s.value
}
