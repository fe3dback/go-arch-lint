package ascii_test

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer/ascii"

	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer/ascii/mocks"
)

//go:embed tests/ascii_input.txt
var testAsciiIn []byte

//go:embed tests/ascii_output.txt
var testAsciiOut []byte

type testModel struct {
	A int    `json:"a"`
	B string `json:"b"`
}

type deps struct {
	asciiColorizer *mocks.MockasciiColorizer
}

type in struct {
	model any
}

func TestASCII_Render(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*deps)
		in      in
		out     string
		wantErr string
	}{
		{
			name: "happy",
			setup: func(d *deps) {
				d.expectCallColor(models.ColorGreen, "1", "green<1>")
				d.expectCallColor(models.ColorRed, "1", "red<1>")
			},
			in:  createIn(),
			out: string(testAsciiOut),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			deps := deps{
				asciiColorizer: mocks.NewMockasciiColorizer(ctrl),
			}
			tt.setup(&deps)

			r := ascii.NewRenderer(
				deps.asciiColorizer,
				map[string]string{
					fmt.Sprintf("%T", tt.in.model): string(testAsciiIn),
				},
			)

			got, gotErr := r.Render(tt.in.model)

			if tt.wantErr != "" {
				require.EqualError(t, gotErr, tt.wantErr)
			} else {
				require.NoError(t, gotErr)
				require.Equal(t, tt.out, got)
			}
		})
	}
}

func createIn(mutators ...func(*in)) in {
	in := in{
		model: createModel(),
	}

	for _, mutate := range mutators {
		mutate(&in)
	}

	return in
}

func createModel() testModel {
	return testModel{
		A: 1,
		B: "hello_world",
	}
}

func (d *deps) expectCallColor(col models.ColorName, text string, want string) {
	d.asciiColorizer.EXPECT().Colorize(col, text).Times(1).Return(want)
}
