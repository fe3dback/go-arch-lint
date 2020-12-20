package spec

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
)

type Provider struct {
	yamlReferenceResolver YamlSourceCodeReferenceResolver
	sourceCode            []byte
}

func NewProvider(
	yamlReferenceResolver YamlSourceCodeReferenceResolver,
	sourceCode []byte,
) *Provider {
	return &Provider{
		yamlReferenceResolver: yamlReferenceResolver,
		sourceCode:            sourceCode,
	}
}

func (sp *Provider) Provide() (arch.Document, error) {
	reader := bytes.NewBuffer(sp.sourceCode)
	decoder := yaml.NewDecoder(
		reader,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(),
		yaml.Strict(),
	)

	document := ArchV1Document{}
	err := decoder.Decode(&document)
	if err != nil {
		return document, fmt.Errorf("can`t parse yaml: %w", err)
	}

	document = document.applyReferences(sp.yamlReferenceResolver)

	// todo: json spec base validation
	//jsonScheme := gojsonschema.NewStringLoader(scheme.V1)
	//jsonDocument := gojsonschema.NewGoLoader(document)
	//
	//result, err := gojsonschema.Validate(jsonScheme, jsonDocument)
	//if err != nil {
	//	return ArchV1{}, fmt.Errorf("failed validate by json scheme: %w", err)
	//}
	//
	//_ = result

	return document, nil
}
