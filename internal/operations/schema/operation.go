package schema

import (
	"encoding/json"
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type Operation struct {
	jsonSchemaProvider jsonSchemaProvider
}

func NewOperation(jsonSchemaProvider jsonSchemaProvider) *Operation {
	return &Operation{
		jsonSchemaProvider: jsonSchemaProvider,
	}
}

func (o *Operation) Behave(in models.CmdSchemaIn) (models.CmdSchemaOut, error) {
	jsonSchema, err := o.jsonSchemaProvider.Provide(in.Version)
	if err != nil {
		return models.CmdSchemaOut{}, fmt.Errorf("failed to provide json schema: %w", err)
	}

	// reformat json to system one line string
	var data interface{}
	err = json.Unmarshal(jsonSchema, &data)
	if err != nil {
		return models.CmdSchemaOut{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	formatted, err := json.Marshal(data)
	if err != nil {
		return models.CmdSchemaOut{}, fmt.Errorf("failed to marshal json: %w", err)
	}

	return models.CmdSchemaOut{
		Version:    in.Version,
		JSONSchema: string(formatted),
	}, nil
}
