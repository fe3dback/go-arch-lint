package renderer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer/mocks"
)

const (
	asciiResult = "ascii:1,2"
	jsonResult  = `{"a":1, "b":2}`
)

type testModel struct {
	A int `json:"a"`
	B int `json:"b"`
}

// ---

type deps struct {
	jsonRenderer  *mocks.MocktypeRenderer
	asciiRenderer *mocks.MocktypeRenderer
}

type in struct {
	outputType models.OutputType
	model      any
}

func TestRenderer_Render(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*deps)
		in      in
		out     string
		wantErr string
	}{
		{
			name: "happy_ascii",
			setup: func(d *deps) {
				d.expectAsciiRendered()
			},
			in:  createIn(models.OutputTypeASCII),
			out: asciiResult,
		},
		{
			name: "happy_json",
			setup: func(d *deps) {
				d.expectJSONRendered()
			},
			in:  createIn(models.OutputTypeJSON),
			out: jsonResult,
		},
		{
			name:    "fail_unknown_output_type",
			setup:   func(d *deps) {},
			in:      createIn("some"),
			wantErr: "unknown renderer type: some",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			deps := deps{
				jsonRenderer:  mocks.NewMocktypeRenderer(ctrl),
				asciiRenderer: mocks.NewMocktypeRenderer(ctrl),
			}
			tt.setup(&deps)

			r := renderer.New(
				deps.jsonRenderer,
				deps.asciiRenderer,
			)

			got, gotErr := r.Render(tt.in.outputType, tt.in.model)

			if tt.wantErr != "" {
				require.EqualError(t, gotErr, tt.wantErr)
			} else {
				require.NoError(t, gotErr)
				require.Equal(t, tt.out, got)
			}
		})
	}
}

func createIn(outType models.OutputType, mutators ...func(*in)) in {
	in := in{
		outputType: outType,
		model:      createModel(),
	}

	for _, mutate := range mutators {
		mutate(&in)
	}

	return in
}

func createModel() testModel {
	return testModel{
		A: 1,
		B: 2,
	}
}

func (d *deps) expectAsciiRendered() {
	d.asciiRenderer.
		EXPECT().
		Render(createModel()).
		Times(1).
		Return(asciiResult, nil)
}

func (d *deps) expectJSONRendered() {
	d.jsonRenderer.
		EXPECT().
		Render(createModel()).
		Times(1).
		Return(jsonResult, nil)
}
