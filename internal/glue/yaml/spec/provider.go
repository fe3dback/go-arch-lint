package spec

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
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

func (sp *Provider) Provide() (arch.Document, []speca.Notice, error) {
	document := ArchV1Document{}

	// read only doc Version
	documentVersion, err := sp.readVersion()
	if err != nil {
		// invalid yaml document
		return document, nil, fmt.Errorf("failed to read 'version' from arch file: %w", err)
	}

	// validate yaml scheme by version
	schemeNotices := sp.jsonSchemeValidate(documentVersion)

	// prepare full document scanner
	reader := bytes.NewBuffer(sp.sourceCode)
	decoder := yaml.NewDecoder(
		reader,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(),
		yaml.Strict(),
	)

	// try to read all document
	err = decoder.Decode(&document)
	if err != nil {
		if len(schemeNotices) > 0 {
			// document invalid, but yaml
			return document, schemeNotices, nil
		}

		// invalid yaml document, or scheme validation failed
		return document, nil, fmt.Errorf("failed to parse arch file (yaml): %w", err)
	}

	document = document.applyReferences(sp.yamlReferenceResolver)
	return document, schemeNotices, nil
}

func (sp *Provider) readVersion() (int, error) {
	type doc struct {
		Version int `yaml:"version"`
	}
	reader := bytes.NewBuffer(sp.sourceCode)
	decoder := yaml.NewDecoder(reader)
	document := doc{}
	err := decoder.Decode(&document)
	if err != nil {
		return 0, err
	}

	return document.Version, nil
}

func (sp *Provider) jsonSchemeValidate(schemeVersion int) []speca.Notice {
	jsonNotices, err := jsonSchemeValidate(sp.sourceCode, schemeVersion)
	if err != nil {
		return []speca.Notice{{
			Notice: fmt.Errorf("failed to validate arch file with json scheme: %v", err),
			Ref:    speca.NewEmptyReference(),
		}}
	}

	schemeNotices := make([]speca.Notice, 0)
	for _, jsonNotice := range jsonNotices {
		schemeRef := speca.NewEmptyReference()
		if jsonNotice.yamlPath != nil {
			schemeRef = sp.yamlReferenceResolver.Resolve(*jsonNotice.yamlPath)
		}

		schemeNotices = append(schemeNotices, speca.Notice{
			Notice: fmt.Errorf(jsonNotice.notice),
			Ref:    schemeRef,
		})
	}

	return schemeNotices
}
