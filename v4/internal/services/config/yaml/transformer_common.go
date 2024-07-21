package yaml

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

func transformFromSyntaxError(tCtx TransformContext, err error) models.Config {
	// todo: why dual FormatError calls?
	ref := extractReferenceFromError(tCtx, err)
	plainErr := yaml.FormatError(err, false, false)

	// bug: formatter (in some cases will print source anyway)
	// try to find source line, ex: ">  9 |   structTags:"
	if strings.Contains(plainErr, "|") && strings.Contains(plainErr, ">") {
		plainErr = strings.Split(plainErr, "\n")[0]
	}

	return models.Config{
		SyntaxProblems: []models.Ref[string]{
			models.NewRef(plainErr, ref),
		},
	}
}

func transform(tCtx TransformContext, doc any) (models.Config, error) {
	switch typed := doc.(type) {
	case ModelV4:
		return transformV4(tCtx, typed), nil
	}

	// todo: add all versions
	return models.Config{}, fmt.Errorf("unknown document version")
}
