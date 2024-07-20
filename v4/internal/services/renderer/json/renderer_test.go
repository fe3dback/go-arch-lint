package json_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer/json"
)

type deps struct {
}

type CmdSomeTestOut struct {
	A int    `json:"a"`
	B string `json:"b"`
}

type in struct {
	model      any
	formatJSON bool
}

func Test_JSON_Render(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*deps)
		in      in
		out     string
		wantErr string
	}{
		{
			name:  "happy_single_line",
			setup: func(d *deps) {},
			in:    createIn(false),
			out:   createOutSingleLine(),
		},
		{
			name:  "happy_formatted",
			setup: func(d *deps) {},
			in:    createIn(true),
			out:   createOutFormatted(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := deps{}
			tt.setup(&deps)

			r := json.NewRenderer()

			got, gotErr := r.Render(tt.in.model, tt.in.formatJSON)

			if tt.wantErr != "" {
				require.EqualError(t, gotErr, tt.wantErr)
			} else {
				require.NoError(t, gotErr)
				require.Equal(t, tt.out, got)
			}
		})
	}
}

func createIn(formatJSON bool, mutators ...func(*in)) in {
	in := in{
		model: CmdSomeTestOut{
			A: 1,
			B: "hello",
		},
		formatJSON: formatJSON,
	}

	for _, mutate := range mutators {
		mutate(&in)
	}

	return in
}

func createOutSingleLine() string {
	return `{"Type":"models.SomeTest","Payload":{"a":1,"b":"hello"}}`
}

func createOutFormatted() string {
	return "{\n  \"Type\": \"models.SomeTest\",\n  \"Payload\": {\n    \"a\": 1,\n    \"b\": \"hello\"\n  }\n}"
}
