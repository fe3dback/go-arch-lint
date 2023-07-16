package spec

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"

	"github.com/goccy/go-yaml"
)

type Provider struct {
	yamlReferenceResolver yamlSourceCodeReferenceResolver
	jsonSchemaProvider    jsonSchemaProvider
}

func NewProvider(
	yamlReferenceResolver yamlSourceCodeReferenceResolver,
	jsonSchemaProvider jsonSchemaProvider,
) *Provider {
	return &Provider{
		yamlReferenceResolver: yamlReferenceResolver,
		jsonSchemaProvider:    jsonSchemaProvider,
	}
}

func (sp *Provider) Provide(archFile string) (arch.Document, []speca.Notice, error) {
	sourceCode, err := os.ReadFile(archFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide source code of archfile: %w", err)
	}

	// read only doc Version
	documentVersion, err := sp.readVersion(sourceCode)
	if err != nil {
		// invalid yaml document
		return nil, nil, fmt.Errorf("failed to read 'version' from arch file: %w", err)
	}

	// validate yaml scheme by version
	schemeNotices := sp.jsonSchemeValidate(documentVersion, sourceCode, archFile)

	// try to read all document
	document, err := sp.decodeDocument(documentVersion, sourceCode, archFile)
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

func (sp *Provider) decodeDocument(version int, sourceCode []byte, filePath string) (arch.Document, error) {
	reader := bytes.NewBuffer(sourceCode)
	decoder := yaml.NewDecoder(
		reader,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(),
		yaml.Strict(),
	)
	resolver := sp.createYamlReferenceResolver(filePath)
	filePathRef := common.NewReferable(filePath, common.NewReferenceSingleLine(filePath, 0, 0))

	// todo: refactor this somehow (dry)
	switch version {
	case 1:
		document := ArchV1Document{filePath: filePathRef}
		err := decoder.Decode(&document)
		if err != nil {
			return nil, err
		}

		return document.applyReferences(resolver), nil
	case 2:
		document := ArchV2Document{filePath: filePathRef}
		err := decoder.Decode(&document)
		if err != nil {
			return nil, err
		}

		return document.applyReferences(resolver), nil
	default:
		document := ArchV3Document{filePath: filePathRef}
		err := decoder.Decode(&document)
		if err != nil {
			return nil, err
		}

		return document.applyReferences(resolver), nil
	}
}

func (sp *Provider) createYamlReferenceResolver(archFilePath string) yamlDocumentPathResolver {
	return func(yamlPath string) common.Reference {
		return sp.yamlReferenceResolver.Resolve(archFilePath, yamlPath)
	}
}

func (sp *Provider) readVersion(sourceCode []byte) (int, error) {
	type doc struct {
		Version int `yaml:"version"`
	}
	reader := bytes.NewBuffer(sourceCode)
	decoder := yaml.NewDecoder(reader)
	document := doc{}
	err := decoder.Decode(&document)
	if err != nil {
		return 0, err
	}

	return document.Version, nil
}

func (sp *Provider) jsonSchemeValidate(schemeVersion int, sourceCode []byte, filePath string) []speca.Notice {
	jsonSchema, err := sp.jsonSchemaProvider.Provide(schemeVersion)
	if err != nil {
		return []speca.Notice{{
			Notice: fmt.Errorf("failed to provide json scheme for validation: %w", err),
			Ref:    common.NewEmptyReference(),
		}}
	}

	jsonNotices, err := jsonSchemeValidate(jsonSchema, sourceCode)
	if err != nil {
		return []speca.Notice{{
			Notice: fmt.Errorf("failed to validate arch file with json scheme: %w", err),
			Ref:    common.NewEmptyReference(),
		}}
	}

	schemeNotices := make([]speca.Notice, 0)
	for _, jsonNotice := range jsonNotices {
		schemeRef := common.NewEmptyReference()
		if jsonNotice.yamlPath != nil {
			schemeRef = sp.yamlReferenceResolver.Resolve(filePath, *jsonNotice.yamlPath)
		}

		schemeNotices = append(schemeNotices, speca.Notice{
			Notice: fmt.Errorf(jsonNotice.notice),
			Ref:    schemeRef,
		})
	}

	return schemeNotices
}
