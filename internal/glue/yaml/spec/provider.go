package spec

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"
)

type Provider struct {
	sourceCode []byte
	validator  Validator
}

func NewProvider(
	sourceCode []byte,
	validator Validator,
) *Provider {
	return &Provider{
		sourceCode: sourceCode,
		validator:  validator,
	}
}

func (sp *Provider) Provide() (ArchDocument, error) {
	reader := bytes.NewBuffer(sp.sourceCode)
	decoder := yaml.NewDecoder(
		reader,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(),
		yaml.Strict(),
	)

	document := Document{}
	err := decoder.Decode(&document)
	if err != nil {
		return ArchDocument{}, fmt.Errorf("can`t parse yaml: %w", err)
	}

	return ArchDocument{
		Document:  document,
		Integrity: sp.validator.Validate(document),
	}, nil
}
