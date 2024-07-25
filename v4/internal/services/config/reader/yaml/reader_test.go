package yaml_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fe3dback/go-arch-lint/v4/internal/services/config/reader/yaml"
)

//go:generate go run ./test_gen/gen_stub.go
// readme: this command will do:
// - read all files in ./test/*.yml
// - parse it with current parser code
// - output parsed DTO to ./test/*_parsed.go
// - this DTO should be checked by human, when changed
// - auto below test will compare latest transformer with stub file

func TestReader_Parse(t *testing.T) {
	tests := []struct {
		testConfig string
		wantError  string
	}{
		{testConfig: "version_below_supported"},
		{testConfig: "version_above_supported"},
		{testConfig: "syntax_problem_elem"},
		{testConfig: "syntax_problem_sys"},
		{testConfig: "3_min"},
		{testConfig: "3_full"},
		{testConfig: "3_deepscan"},
		{testConfig: "4_min"},
		{testConfig: "4_full"},
		{testConfig: "4_tags_list"},
	}
	for _, tt := range tests {
		reader := yaml.NewReader()

		t.Run(tt.testConfig, func(t *testing.T) {
			sourceCode, err := os.ReadFile(fmt.Sprintf("./test/%s.yml", tt.testConfig))
			require.NoError(t, err)

			conf, err := reader.Parse("/conf.yml", sourceCode)

			if tt.wantError != "" {
				require.EqualError(t, err, tt.wantError)
			} else {
				require.NoError(t, err)
			}

			encoded, err := json.MarshalIndent(conf, "", "  ")
			require.NoError(t, err)

			savedStub, err := os.ReadFile(fmt.Sprintf("./test/%s_parsed.json", tt.testConfig))
			require.NoError(t, err)

			require.Equal(t, savedStub, encoded)
		})
	}
}
