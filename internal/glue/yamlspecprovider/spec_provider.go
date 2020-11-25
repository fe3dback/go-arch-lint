package yamlspecprovider

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"
)

type YamlSpecProvider struct {
	sourceCode []byte
}

func NewYamlSpecProvider(sourceCode []byte) *YamlSpecProvider {
	return &YamlSpecProvider{
		sourceCode: sourceCode,
	}
}

func (sp *YamlSpecProvider) Provide() (*YamlSpec, error) {
	reader := bytes.NewBuffer(sp.sourceCode)
	decoder := yaml.NewDecoder(
		reader,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(),
		yaml.Strict(),
	)

	spec := YamlSpec{}
	err := decoder.Decode(&spec)
	if err != nil {
		return nil, fmt.Errorf("can`t parse yaml: %w", err)
	}

	return &spec, nil
}
