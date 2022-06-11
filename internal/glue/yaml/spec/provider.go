package spec

import (
	"bytes"
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"

	"github.com/goccy/go-yaml"
)

type Provider struct {
	yamlReferenceResolver YAMLSourceCodeReferenceResolver
	jsonSchemaProvider    JSONSchemaProvider
	sourceCode            []byte
}

func NewProvider(
	yamlReferenceResolver YAMLSourceCodeReferenceResolver,
	jsonSchemaProvider JSONSchemaProvider,
	sourceCode []byte,
) *Provider {
	return &Provider{
		yamlReferenceResolver: yamlReferenceResolver,
		jsonSchemaProvider:    jsonSchemaProvider,
		sourceCode:            sourceCode,
	}
}

func (sp *Provider) Provide() (arch.Document, []speca.Notice, error) {
	// read only doc Version
	documentVersion, err := sp.readVersion()
	if err != nil {
		// invalid yaml document
		return nil, nil, fmt.Errorf("failed to read 'version' from arch file: %w", err)
	}

	// validate yaml scheme by version
	schemeNotices := sp.jsonSchemeValidate(documentVersion)

	// try to read all document
	document, err := sp.decodeDocument(documentVersion)
	if err != nil {
		if len(schemeNotices) > 0 {
			// document invalid, but yaml
			return document, schemeNotices, nil
		}

		// invalid yaml document, or scheme validation failed
		return nil, nil, fmt.Errorf("failed to parse arch file (yaml): %w", err)
	}

	return document, schemeNotices, nil
}

func (sp *Provider) decodeDocument(version int) (arch.Document, error) {
	reader := bytes.NewBuffer(sp.sourceCode)
	decoder := yaml.NewDecoder(
		reader,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(),
		yaml.Strict(),
	)

	// todo: refactor this somehow (dry)
	switch version {
	case 1:
		document := ArchV1Document{}
		err := decoder.Decode(&document)
		if err != nil {
			return nil, err
		}

		return document.applyReferences(sp.yamlReferenceResolver), nil
	case 2:
		document := ArchV2Document{}
		err := decoder.Decode(&document)
		if err != nil {
			return nil, err
		}

		return document.applyReferences(sp.yamlReferenceResolver), nil
	default:
		document := ArchV3Document{}
		err := decoder.Decode(&document)
		if err != nil {
			return nil, err
		}

		return document.applyReferences(sp.yamlReferenceResolver), nil
	}
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
	jsonSchema, err := sp.jsonSchemaProvider.Provide(schemeVersion)
	if err != nil {
		return []speca.Notice{{
			Notice: fmt.Errorf("failed to provide json scheme for validation: %w", err),
			Ref:    speca.NewEmptyReference(),
		}}
	}

	jsonNotices, err := jsonSchemeValidate(jsonSchema, sp.sourceCode)
	if err != nil {
		return []speca.Notice{{
			Notice: fmt.Errorf("failed to validate arch file with json scheme: %w", err),
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
