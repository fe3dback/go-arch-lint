package ascii

import (
	"bytes"
	"fmt"
	"text/template"

	"golang.org/x/exp/maps"
)

const (
	fnColorize   = "colorize"
	fnTrimPrefix = "trimPrefix"
	fnTrimSuffix = "trimSuffix"
	fnTrimDef    = "def"
	fnPadLeft    = "padLeft"
	fnPadRight   = "padRight"
	fnLinePrefix = "linePrefix"
	fnDir        = "dir"
	fnPlus       = "plus"
	fnMinus      = "minus"
	fnConcat     = "concat"
)

type ASCII struct {
	asciiColorizer asciiColorizer
	templates      map[string]string
}

func NewRenderer(
	asciiColorizer asciiColorizer,
	templates map[string]string,
) *ASCII {
	return &ASCII{
		asciiColorizer: asciiColorizer,
		templates:      templates,
	}
}

func (r *ASCII) Render(model any) (string, error) {
	templateName := fmt.Sprintf("%T", model)
	templateBuffer, exist := r.templates[templateName]

	if !exist {
		return "", fmt.Errorf("ascii template for model '%s' not exist. Found:[%#v]", templateName, maps.Keys(r.templates))
	}

	tpl, err := template.
		New(templateName).
		Funcs(map[string]interface{}{
			fnColorize:   r.asciiColorize,
			fnTrimPrefix: r.asciiTrimPrefix,
			fnTrimSuffix: r.asciiTrimSuffix,
			fnTrimDef:    r.asciiDefaultValue,
			fnPadLeft:    r.asciiPadLeft,
			fnPadRight:   r.asciiPadRight,
			fnLinePrefix: r.asciiLinePrefix,
			fnDir:        r.asciiPathDirectory,
			fnPlus:       r.asciiPlus,
			fnMinus:      r.asciiMinus,
			fnConcat:     r.asciiConcat,
		}).
		Parse(
			preprocessRawASCIITemplate(templateBuffer),
		)
	if err != nil {
		return "", fmt.Errorf("failed to render ascii view '%s': %w", templateName, err)
	}

	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, model)
	if err != nil {
		return "", fmt.Errorf("failed to execute template '%s': %w", templateName, err)
	}

	return buffer.String(), nil
}
