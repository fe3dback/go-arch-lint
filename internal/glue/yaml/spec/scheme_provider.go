package spec

import (
	"github.com/xeipuuv/gojsonschema"

	"github.com/fe3dback/go-arch-lint/internal/scheme"
)

func provideScheme(version int) *gojsonschema.JSONLoader {
	var jsonScheme string

	switch version {
	case 1:
		jsonScheme = scheme.V1
	default:
		return nil
	}

	loader := gojsonschema.NewStringLoader(jsonScheme)
	return &loader
}
