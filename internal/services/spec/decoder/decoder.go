package decoder

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
	"github.com/fe3dback/go-yaml"
)

type Decoder struct {
	yamlReferenceResolver yamlSourceCodeReferenceResolver
	jsonSchemaProvider    jsonSchemaProvider
}

func NewDecoder(
	yamlReferenceResolver yamlSourceCodeReferenceResolver,
	jsonSchemaProvider jsonSchemaProvider,
) *Decoder {
	return &Decoder{
		yamlReferenceResolver: yamlReferenceResolver,
		jsonSchemaProvider:    jsonSchemaProvider,
	}
}

func (sp *Decoder) Decode(archFile string) (spec.Document, []arch.Notice, error) {
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

func (sp *Decoder) decodeDocument(version int, sourceCode []byte, filePath string) (spec.Document, error) {
	reader := bytes.NewBuffer(sourceCode)
	decoder := yaml.NewDecoder(
		reader,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(),
		yaml.Strict(),
	)

	decodeCtx := context.WithValue(context.Background(), yamlParentFileCtx{}, filePath)
	document := sp.createEmptyDocumentBeVersion(version)

	err := decoder.DecodeContext(decodeCtx, document)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (sp *Decoder) createEmptyDocumentBeVersion(version int) spec.Document {
	switch version {
	case 1:
		return &ArchV1{}
	case 2:
		return &ArchV2{}
	}

	// latest be default (it will be rejected next in spec validator, if version is not v3)
	return &ArchV3{}
}

func (sp *Decoder) readVersion(sourceCode []byte) (int, error) {
	type doc struct {
		Version int `json:"version"`
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

func (sp *Decoder) jsonSchemeValidate(schemeVersion int, sourceCode []byte, filePath string) []arch.Notice {
	jsonSchema, err := sp.jsonSchemaProvider.Provide(schemeVersion)
	if err != nil {
		return []arch.Notice{{
			Notice: fmt.Errorf("failed to provide json scheme for validation: %w", err),
			Ref:    common.NewEmptyReference(),
		}}
	}

	jsonNotices, err := jsonSchemeValidate(jsonSchema, sourceCode)
	if err != nil {
		return []arch.Notice{{
			Notice: fmt.Errorf("failed to validate arch file with json scheme: %w", err),
			Ref:    common.NewEmptyReference(),
		}}
	}

	schemeNotices := make([]arch.Notice, 0)
	for _, jsonNotice := range jsonNotices {
		schemeRef := common.NewEmptyReference()
		if jsonNotice.yamlPath != nil {
			schemeRef = sp.yamlReferenceResolver.Resolve(filePath, *jsonNotice.yamlPath)
		}

		schemeNotices = append(schemeNotices, arch.Notice{
			Notice: fmt.Errorf(jsonNotice.notice),
			Ref:    schemeRef,
		})
	}

	return schemeNotices
}
