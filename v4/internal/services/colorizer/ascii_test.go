package colorizer_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/colorizer"
)

const (
	sourceText     = "hello"
	sourceTextRed  = "\x1b[91mhello\x1b[0m"
	sourceTextBlue = "\x1b[94mhello\x1b[0m"
)

type deps struct {
}

type in struct {
	useColors bool
	color     models.ColorName
	text      string
}

func TestASCII_Colorize(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*deps)
		in    in
		out   string
	}{
		{
			name: "happy_no_colors",
			setup: func(d *deps) {
			},
			in: createIn(models.ColorRed, func(in *in) {
				in.useColors = false
			}),
			out: sourceText,
		},
		{
			name: "happy_red",
			setup: func(d *deps) {
			},
			in:  createIn(models.ColorRed),
			out: sourceTextRed,
		},
		{
			name: "happy_blue",
			setup: func(d *deps) {
			},
			in:  createIn(models.ColorBlue),
			out: sourceTextBlue,
		},
		{
			name:  "unknown_color",
			setup: func(d *deps) {},
			in:    createIn("not-exist-color"),
			out:   sourceText,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := deps{}
			tt.setup(&deps)

			r := colorizer.New(tt.in.useColors)

			got := r.Colorize(tt.in.color, tt.in.text)
			require.Equal(t, tt.out, got)
		})
	}
}

func createIn(col models.ColorName, mutators ...func(*in)) in {
	in := in{
		color:     col,
		text:      sourceText,
		useColors: true,
	}

	for _, mutate := range mutators {
		mutate(&in)
	}

	return in
}
