package yaml

import (
	"bytes"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) ReadFile(path models.PathAbsolute) (models.Config, error) {
	sourceCode, err := os.ReadFile(string(path))
	if err != nil {
		return models.Config{}, fmt.Errorf("failed to provide source code of archfile: %w", err)
	}

	return r.Parse(path, sourceCode)
}

func (r *Reader) Parse(path models.PathAbsolute, source []byte) (models.Config, error) {
	tCtx := TransformContext{
		file:   path,
		source: source,
	}

	// read only doc Version
	documentVersion, err := r.readVersion(source)
	if err != nil {
		return transformFromSyntaxError(tCtx, fmt.Errorf("failed to read 'version' from arch file: %w", err)), nil
	}

	// try to read all document
	document, err := r.decodeDocument(documentVersion, source)
	if err != nil {
		return transformFromSyntaxError(tCtx, err), nil
	}

	return transform(tCtx, document)
}

func (r *Reader) readVersion(sourceCode []byte) (int, error) {
	type doc struct {
		Version int `json:"version"`
	}
	reader := bytes.NewBuffer(sourceCode)
	decoder := yaml.NewDecoder(reader)
	document := doc{}

	err := decoder.Decode(&document)
	if err != nil {
		return 0, fmt.Errorf("failed decode version from yaml doc: %w", err)
	}

	return document.Version, nil
}

func (r *Reader) decodeDocument(version int, sourceCode []byte) (any, error) {
	reader := bytes.NewBuffer(sourceCode)
	decoder := yaml.NewDecoder(
		reader,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(),
		yaml.Strict(),
	)

	document := r.createEmptyDocumentBeVersion(version)

	err := decoder.Decode(document)
	if err != nil {
		return nil, fmt.Errorf("failed decode yaml: %w", err)
	}

	return document, nil
}

func (r *Reader) createEmptyDocumentBeVersion(version int) any {
	switch version {
	case 1:
		// todo
		panic("not implemented v1")
	case 2:
		// todo
		panic("not implemented v2")
	case 3:
		// todo
		panic("not implemented v3")
	}

	// latest be default (it will be rejected next in spec validator, if version is not v4)
	return &ModelV4{}
}
