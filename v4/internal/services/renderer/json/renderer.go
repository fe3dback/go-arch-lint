package json

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JSON struct {
}

func NewRenderer() *JSON {
	return &JSON{}
}

func (r *JSON) Render(model any, formatJSON bool) (string, error) {
	var jsonBuffer []byte
	var marshalErr error

	modelType, err := r.extractModelType(model)
	if err != nil {
		return "", fmt.Errorf("failed extract model type from '%T' (maybe not matched pattern: 'CmdXXXOut') : %w", model, err)
	}

	wrapperModel := struct {
		Type    string      `json:"Type"`
		Payload interface{} `json:"Payload"`
	}{
		Type:    modelType,
		Payload: model,
	}

	if !formatJSON {
		jsonBuffer, marshalErr = json.Marshal(wrapperModel)
	} else {
		jsonBuffer, marshalErr = json.MarshalIndent(wrapperModel, "", "  ")
	}

	if marshalErr != nil {
		return "", fmt.Errorf("failed to marshal payload '%v' to json: %w", model, marshalErr)
	}

	return string(jsonBuffer), nil
}

// Rename "anypackage.CmdXXXXOut" to "models.XXXX"
// for back compatible with previous response version
func (r *JSON) extractModelType(model any) (string, error) {
	const expectedPrefix = "Cmd"
	const expectedSuffix = "Out"

	alias := fmt.Sprintf("%T", model)
	dotIndex := strings.Index(alias, ".")

	if dotIndex == -1 {
		return "", fmt.Errorf("DTO type '%s' without package name", alias)
	}

	dtoName := alias[dotIndex+1:]

	if !strings.HasPrefix(dtoName, expectedPrefix) {
		return "", fmt.Errorf("DTO name '%s' alias '%s' should has prefix '%s'", dtoName, alias, expectedPrefix)
	}

	if !strings.HasSuffix(dtoName, expectedSuffix) {
		return "", fmt.Errorf("DTO name '%s' alias '%s' should has suffix '%s'", dtoName, alias, expectedSuffix)
	}

	return fmt.Sprintf(
		"models.%s",
		strings.TrimPrefix(strings.TrimSuffix(dtoName, expectedSuffix), expectedPrefix),
	), nil
}
