package spec

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/xeipuuv/gojsonschema"

	"github.com/fe3dback/go-arch-lint/internal/models"
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
		return document, nil, fmt.Errorf("can`t parse yaml: %w", err)
	}

	// validate yaml scheme by version
	schemeNotices := sp.validateJsonScheme(documentVersion)

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
		return document, nil, fmt.Errorf("can`t parse yaml: %w", err)
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

func (sp *Provider) validateJsonScheme(version int) []speca.Notice {
	var body interface{}
	err := yaml.Unmarshal(sp.sourceCode, &body)
	if err != nil {
		// invalid yaml document
		return nil
	}

	jsonScheme := provideScheme(version)
	if jsonScheme == nil {
		// unknown spec version, skip scheme check
		return nil
	}

	jsonBody, err := json.Marshal(&body)
	if err != nil {
		// invalid json struct in mem
		return nil
	}

	jsonDocument := gojsonschema.NewBytesLoader(jsonBody)
	result, err := gojsonschema.Validate(*jsonScheme, jsonDocument)
	if err != nil {
		return []speca.Notice{{
			Notice: fmt.Errorf("failed to validate arch file with json scheme: %v", err),
			Ref:    speca.NewEmptyReference(),
		}}
	}

	notices := make([]speca.Notice, 0)
	for _, err := range result.Errors() {
		notices = append(notices, sp.jsonSchemeErrorToNotice(err))
	}

	return notices
}

func (sp *Provider) jsonSchemeErrorToNotice(err gojsonschema.ResultError) speca.Notice {
	return speca.Notice{
		Notice: fmt.Errorf(err.String()),
		Ref:    sp.referenceByJsonSchemeError(err),
	}
}

func (sp *Provider) referenceByJsonSchemeError(err gojsonschema.ResultError) models.Reference {
	// todo: check map and slice path's

	// root
	path := "$"

	// context
	if err.Field() == "(root)" {
		propertyName, ok := err.Details()["property"]
		if !ok {
			return speca.NewEmptyReference()
		}

		path = fmt.Sprintf("%s.%s", path, propertyName)
	} else {
		path = fmt.Sprintf("%s.%s", path, err.Field())
	}

	// resolve path
	return sp.yamlReferenceResolver.Resolve(path)
}
