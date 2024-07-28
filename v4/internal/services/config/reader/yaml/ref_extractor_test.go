package yaml

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

const confFile = "/some/path/conf.yml"

func Test_extractReferenceFromError(t *testing.T) {
	tests := []struct {
		text string
		want models.Reference
	}{
		{
			text: `[24:26] unknown field "in2"
  21 |     - "^.*\/test\/.*$"
  22 | 
  23 | vendors:
> 24 |   go-common:           { in2: golang.org/x/sync/errgroup }
                                ^
  25 |   go-ast:              { in: [ golang.org/x/mod/modfile, golang.org/x/tools/go/packages ] }
  26 |   3rd-cobra:           { in: github.com/spf13/cobra }
  27 |   3rd-color-fmt:       { in: github.com/logrusorgru/aurora/v3 }
`,
			want: models.Reference{
				File:   confFile,
				Line:   24,
				Column: 26,
				Valid:  true,
			},
		},
		{
			text: `[22:1] cannot unmarshal map[string]interface {} into Go struct field ModelV4.Exclude of type string
  19 |   files:
  20 |     - "^.*_test\\.go$"
  21 |     - "^.*\/test\/.*$":
       ^
  23 | vendors:
  24 |   go-common:           { in: golang.org/x/sync/errgroup }
  25 |   go-ast:              { in: [ golang.org/x/mod/modfile, golang.org/x/tools/go/packages ] }`,
			want: models.Reference{
				File:   confFile,
				Line:   22,
				Column: 1,
				Valid:  true,
			},
		},
	}
	for ind, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", ind), func(t *testing.T) {
			lines := strings.Split(tt.text, "\n")
			fCtx := TransformContext{
				file: confFile,
			}

			got := extractReferenceFromError(fCtx, lines[1:])
			require.Equal(t, tt.want, got)
		})
	}
}
