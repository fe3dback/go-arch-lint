package yaml

import (
	"bytes"

	"github.com/goccy/go-yaml"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type TransformContext struct {
	file   models.PathAbsolute
	source []byte
}

func (tc *TransformContext) createReference(documentAstPath string) models.Reference {
	astPath, err := yaml.PathString(documentAstPath)
	if err != nil {
		return models.NewInvalidReference()
	}

	astNode, err := astPath.ReadNode(bytes.NewReader(tc.source))
	if err != nil {
		return models.NewInvalidReference()
	}

	tok := astNode.GetToken()
	return models.NewReference(
		tc.file,
		tok.Position.Line,
		tok.Position.Column,
	)
}
