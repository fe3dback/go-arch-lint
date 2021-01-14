package printer

import (
	"github.com/logrusorgru/aurora/v3"
)

type ColorPrinter struct {
	au aurora.Aurora
}

func NewColorPrinter(au aurora.Aurora) *ColorPrinter {
	return &ColorPrinter{
		au: au,
	}
}

func (cp *ColorPrinter) Red(in string) (out string) {
	return cp.au.Red(in).String()
}

func (cp *ColorPrinter) Green(in string) (out string) {
	return cp.au.Green(in).String()
}

func (cp *ColorPrinter) Yellow(in string) (out string) {
	return cp.au.Yellow(in).String()
}

func (cp *ColorPrinter) Blue(in string) (out string) {
	return cp.au.Blue(in).String()
}

func (cp *ColorPrinter) Magenta(in string) (out string) {
	return cp.au.Magenta(in).String()
}

func (cp *ColorPrinter) Cyan(in string) (out string) {
	return cp.au.Cyan(in).String()
}

func (cp *ColorPrinter) White(in string) (out string) {
	return cp.au.White(in).String()
}

func (cp *ColorPrinter) Gray(in string) (out string) {
	return cp.au.BrightBlack(in).String()
}
