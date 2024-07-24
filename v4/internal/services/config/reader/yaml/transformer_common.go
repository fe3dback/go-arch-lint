package yaml

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

func transformFromSyntaxError(tCtx TransformContext, err error) models.Config {
	errText := yaml.FormatError(err, false, true)
	errLines := strings.Split(errText, "\n")
	plainErr := errLines[0]

	ref := extractReferenceFromError(tCtx, errLines[1:])

	return models.Config{
		SyntaxProblems: []models.Ref[string]{
			models.NewRef(plainErr, ref),
		},
	}
}

func transform(tCtx TransformContext, doc any) (models.Config, error) {
	switch typed := doc.(type) {
	case *ModelV4:
		return transformV4(tCtx, *typed), nil
	}

	// todo: add all versions
	return models.Config{}, fmt.Errorf("unknown document type %T", doc)
}
