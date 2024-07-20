package colorizer_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/colorizer"
)

const (
	sourceText     = "hello"
	sourceTextRed  = "hello"
	sourceTextBlue = "hello"
)

type deps struct {
}

type in struct {
	color models.ColorName
	text  string
}

func TestASCII_Colorize(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*deps)
		in    in
		out   string
	}{
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

			r := colorizer.New()

			got := r.Colorize(tt.in.color, tt.in.text)
			require.Equal(t, []byte(tt.out), []byte(got))
		})
	}
}

func createIn(col models.ColorName, mutators ...func(*in)) in {
	in := in{
		color: col,
		text:  sourceText,
	}

	for _, mutate := range mutators {
		mutate(&in)
	}

	return in
}
