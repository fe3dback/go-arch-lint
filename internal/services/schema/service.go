package schema

import (
	"encoding/json"
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Behave(schema models.FlagsSchema) (models.Schema, error) {
	// reformat json to system one line string
	var data interface{}
	err := json.Unmarshal([]byte(schema.JSONSchema), &data)
	if err != nil {
		return models.Schema{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	formatted, err := json.Marshal(data)
	if err != nil {
		return models.Schema{}, fmt.Errorf("failed to marshal json: %w", err)
	}

	return models.Schema{
		Version:    schema.Version,
		JSONSchema: string(formatted),
	}, nil
}
