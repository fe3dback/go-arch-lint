package colorizer

import (
	"os"

	"github.com/muesli/termenv"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

var (
	env          *termenv.Output
	colorRed     termenv.Color
	colorGreen   termenv.Color
	colorBlue    termenv.Color
	colorYellow  termenv.Color
	colorMagenta termenv.Color
	colorCyan    termenv.Color
	colorGray    termenv.Color
)

func init() {
	env = termenv.NewOutput(os.Stdout)
	colorRed = env.Color("#d62d20")
	colorGreen = env.Color("#008744")
	colorBlue = env.Color("#0057e7")
	colorYellow = env.Color("#ffa700")
	colorMagenta = env.Color("#d62976")
	colorCyan = env.Color("#2691a5")
	colorGray = env.Color("#777777")
}

type ASCII struct {
	useColors bool
}

func New(useColors bool) *ASCII {
	return &ASCII{
		useColors: useColors,
	}
}

func (c *ASCII) Colorize(color models.ColorName, text string) string {
	if !c.useColors {
		return text
	}

	s := env.String(text)

	switch color {
	case models.ColorRed:
		return s.Foreground(colorRed).String()
	case models.ColorGreen:
		return s.Foreground(colorGreen).String()
	case models.ColorYellow:
		return s.Foreground(colorYellow).String()
	case models.ColorBlue:
		return s.Foreground(colorBlue).String()
	case models.ColorMagenta:
		return s.Foreground(colorMagenta).String()
	case models.ColorCyan:
		return s.Foreground(colorCyan).String()
	case models.ColorGray:
		return s.Foreground(colorGray).String()
	default:
		return text
	}
}
