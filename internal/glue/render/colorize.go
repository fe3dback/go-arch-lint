package render

import (
	"fmt"
)

const (
	colorRed     colorName = "red"
	colorGreen   colorName = "green"
	colorYellow  colorName = "yellow"
	colorBlue    colorName = "blue"
	colorMagenta colorName = "magenta"
	colorCyan    colorName = "cyan"
	colorGray    colorName = "gray"
)

type (
	colorizer struct {
		printer ColorPrinter
	}

	colorName = string
)

func newColorizer(printer ColorPrinter) *colorizer {
	return &colorizer{
		printer: printer,
	}
}

func (c *colorizer) colorize(color colorName, input string) (string, error) {
	switch color {
	case colorRed:
		return c.printer.Red(input), nil
	case colorGreen:
		return c.printer.Green(input), nil
	case colorYellow:
		return c.printer.Yellow(input), nil
	case colorBlue:
		return c.printer.Blue(input), nil
	case colorMagenta:
		return c.printer.Magenta(input), nil
	case colorCyan:
		return c.printer.Cyan(input), nil
	case colorGray:
		return c.printer.Gray(input), nil
	default:
		return "", fmt.Errorf("invalid color '%s'", color)
	}
}
