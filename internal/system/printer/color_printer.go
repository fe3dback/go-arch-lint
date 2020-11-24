package printer

import (
	"fmt"

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
	return fmt.Sprintf("%s", cp.au.Red(in))
}

func (cp *ColorPrinter) Green(in string) (out string) {
	return fmt.Sprintf("%s", cp.au.Green(in))
}

func (cp *ColorPrinter) Yellow(in string) (out string) {
	return fmt.Sprintf("%s", cp.au.Yellow(in))
}

func (cp *ColorPrinter) Blue(in string) (out string) {
	return fmt.Sprintf("%s", cp.au.Blue(in))
}

func (cp *ColorPrinter) Magenta(in string) (out string) {
	return fmt.Sprintf("%s", cp.au.Magenta(in))
}

func (cp *ColorPrinter) Cyan(in string) (out string) {
	return fmt.Sprintf("%s", cp.au.Cyan(in))
}

func (cp *ColorPrinter) White(in string) (out string) {
	return fmt.Sprintf("%s", cp.au.White(in))
}
