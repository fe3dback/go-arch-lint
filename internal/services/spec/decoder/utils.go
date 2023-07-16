package decoder

import (
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
)

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
