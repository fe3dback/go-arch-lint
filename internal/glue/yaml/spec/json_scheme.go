package spec

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	"github.com/fe3dback/go-arch-lint/internal/scheme"
)

type (
	jsonSchemeNotice struct {
		notice   string
		yamlPath *string
	}
)

func jsonSchemeValidate(sourceCode []byte, version int) ([]jsonSchemeNotice, error) {
	jsonScheme, err := jsonSchemeByVersion(version)
	if err != nil {
		return nil, fmt.Errorf("failed provide json scheme loader: %w", err)
	}

	jsonDocument, err := jsonSchemeDocumentByCode(sourceCode)
	if err != nil {
		return nil, fmt.Errorf("failed provide json document loader: %w", err)
	}

	if jsonScheme == nil || jsonDocument == nil {
		// unknown spec version, skip scheme check
		return nil, fmt.Errorf("json document or scheme is invalid")
	}

	result, err := gojsonschema.Validate(*jsonScheme, *jsonDocument)
	if err != nil {
		return nil, fmt.Errorf("json scheme validation error: %w", err)
	}

	notices := make([]jsonSchemeNotice, 0)
	for _, schemeErr := range result.Errors() {
		yamlPath := jsonSchemeExtractYamlPathFromError(schemeErr)
		titlePath := fmt.Sprintf("? <%s>", schemeErr.Context().String())

		if yamlPath != nil {
			titlePath = *yamlPath
		}

		notices = append(notices, jsonSchemeNotice{
			notice:   fmt.Sprintf("(%s) %s", titlePath, schemeErr.Description()),
			yamlPath: yamlPath,
		})
	}

	return notices, nil
}

func jsonSchemeByVersion(version int) (*gojsonschema.JSONLoader, error) {
	var jsonScheme string

	switch version {
	case 1:
		jsonScheme = scheme.V1
	case 2:
		jsonScheme = scheme.V2
	default:
		return nil, fmt.Errorf("unknown version: %d", version)
	}

	loader := gojsonschema.NewStringLoader(jsonScheme)
	return &loader, nil
}

func jsonSchemeDocumentByCode(sourceCode []byte) (*gojsonschema.JSONLoader, error) {
	jsonBody, err := yamlToJson(sourceCode)
	if err != nil {
		return nil, fmt.Errorf("failed transform yaml to json: %w", err)
	}

	loader := gojsonschema.NewBytesLoader(jsonBody)
	return &loader, nil
}

func jsonSchemeExtractYamlPathFromError(err gojsonschema.ResultError) *string {
	// todo: map's path not working, because json path $.a.b.c in yaml, can be:
	// - $.a.b.c (object)
	// - $.a[b].c (map)

	// root
	path := "(root)"

	// context
	if err.Field() == "(root)" {
		propertyName, ok := err.Details()["property"]
		if !ok {
			return nil
		}

		path = fmt.Sprintf("%s.%s", path, propertyName)
	} else {
		path = fmt.Sprintf("%s.%s", path, err.Field())
	}

	path = jsonSchemeTransformJsonPathToYamlPath(path)

	// resolve path
	return &path
}

// transform jsonPath to yamlPath
//   "(root).exclude.1" 		-> "$.exclude[1]"
//   "(root).some.field.22" 	-> "$.some.field[22]"
//   "(root).some.field.22a.b" 	-> "$.some.field.22a.b"
func jsonSchemeTransformJsonPathToYamlPath(path string) string {
	// root -> $
	path = strings.Replace(path, "(root)", "$", 1)

	// array index .1 -> [1]
	re := regexp.MustCompile(`\.([0-9]+)(\.|$)`)
	path = re.ReplaceAllString(path, "[${1}]${2}")

	return path
}
