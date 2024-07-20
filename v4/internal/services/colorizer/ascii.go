package colorizer

import (
	cl "github.com/fatih/color"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type ASCII struct {
}

func New() *ASCII {
	return &ASCII{}
}

func (c *ASCII) Colorize(color models.ColorName, text string) string {
	switch color {
	case models.ColorRed:
		return cl.HiRedString(text)
	case models.ColorGreen:
		return cl.HiGreenString(text)
	case models.ColorYellow:
		return cl.HiYellowString(text)
	case models.ColorBlue:
		return cl.HiBlueString(text)
	case models.ColorMagenta:
		return cl.HiMagentaString(text)
	case models.ColorCyan:
		return cl.HiCyanString(text)
	case models.ColorGray:
		return cl.HiWhiteString(text)
	default:
		return text
	}
}
