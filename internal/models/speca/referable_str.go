package speca

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/goccy/go-yaml"
)

type (
	ReferableString struct {
		value string
		ref   models.Reference
	}
)

func (s *ReferableString) UnmarshalYAML(bytes []byte) error {
	return yaml.Unmarshal(bytes, &s.value)
}

func (s *ReferableString) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(s.value)
}

func NewReferableString(value string, ref models.Reference) ReferableString {
	return ReferableString{value: value, ref: ref}
}

func NewEmptyReferableString(value string) ReferableString {
	return ReferableString{value: value, ref: NewEmptyReference()}
}

func (s ReferableString) Reference() models.Reference {
	return s.ref
}

func (s ReferableString) Value() string {
	return s.value
}
