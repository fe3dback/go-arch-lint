package speca

import (
	"github.com/goccy/go-yaml"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	ReferableInt struct {
		value int
		ref   models.Reference
	}
)

func (s ReferableInt) UnmarshalYAML(bytes []byte) error {
	return yaml.Unmarshal(bytes, &s.value)
}

func (s ReferableInt) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(s.value)
}

func NewReferableInt(value int, ref models.Reference) ReferableInt {
	return ReferableInt{value: value, ref: ref}
}

func NewEmptyReferableInt(value int) ReferableInt {
	return ReferableInt{value: value, ref: NewEmptyReference()}
}

func (s ReferableInt) Reference() models.Reference {
	return s.ref
}

func (s ReferableInt) Value() int {
	return s.value
}
