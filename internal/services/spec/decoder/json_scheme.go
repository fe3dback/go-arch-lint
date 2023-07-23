package decoder

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/fe3dback/go-yaml"
	"github.com/xeipuuv/gojsonschema"
)

type (
	jsonSchemeNotice struct {
		notice   string
		yamlPath *string
	}
)

func jsonSchemeValidate(jsonSchema []byte, sourceCode []byte) ([]jsonSchemeNotice, error) {
	jsonDocument, err := jsonSchemeDocumentByCode(sourceCode)
	if err != nil {
		return nil, fmt.Errorf("failed provide json document loader: %w", err)
	}
	if jsonDocument == nil {
		return nil, fmt.Errorf("json document is invalid")
	}

	jsonScheme := gojsonschema.NewBytesLoader(jsonSchema)
	result, err := gojsonschema.Validate(jsonScheme, *jsonDocument)
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

func jsonSchemeDocumentByCode(sourceCode []byte) (*gojsonschema.JSONLoader, error) {
	jsonBody, err := yamlToJSON(sourceCode)
	if err != nil {
		return nil, fmt.Errorf("failed transform yaml to json: %w", err)
	}

	loader := gojsonschema.NewBytesLoader(jsonBody)
	return &loader, nil
}

func yamlToJSON(sourceCode []byte) ([]byte, error) {
	var body interface{}
	err := yaml.Unmarshal(sourceCode, &body)
	if err != nil {
		// invalid yaml document
		return nil, fmt.Errorf("invalid source yaml: %w", err)
	}

	jsonBody, err := json.Marshal(&body)
	if err != nil {
		// invalid json struct in mem
		return nil, fmt.Errorf("failed marshal to json: %w", err)
	}

	return jsonBody, nil
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

	path = jsonSchemeTransformJSONPathToYamlPath(path)

	// resolve path
	return &path
}

// transform jsonPath to yamlPath
//
//	"(root).exclude.1" 		-> "$.exclude[1]"
//	"(root).some.field.22" 	-> "$.some.field[22]"
//	"(root).some.field.22a.b" 	-> "$.some.field.22a.b"
func jsonSchemeTransformJSONPathToYamlPath(path string) string {
	// root -> $
	path = strings.Replace(path, "(root)", "$", 1)

	// array index .1 -> [1]
	re := regexp.MustCompile(`\.([0-9]+)(\.|$)`)
	path = re.ReplaceAllString(path, "[${1}]${2}")

	return path
}
