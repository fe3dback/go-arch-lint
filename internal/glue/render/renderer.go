package render

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"text/template"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

const (
	fnColorize   = "colorize"
	fnTrimPrefix = "trimPrefix"
	fnTrimSuffix = "trimSuffix"
	fnTrimDef    = "def"
	fnPadLeft    = "padLeft"
	fnPadRight   = "padRight"
	fnDir        = "dir"
	fnPlus       = "plus"
	fnMinus      = "minus"
)

type (
	Renderer struct {
		colorPrinter      ColorPrinter
		outputType        models.OutputType
		outputJSONOneLine bool
		asciiTemplates    map[string]string
	}
)

func NewRenderer(
	colorPrinter ColorPrinter,
	outputType models.OutputType,
	outputJSONOneLine bool,
	asciiTemplates map[string]string,
) *Renderer {
	return &Renderer{
		colorPrinter:      colorPrinter,
		outputType:        outputType,
		outputJSONOneLine: outputJSONOneLine,
		asciiTemplates:    asciiTemplates,
	}
}

func (r *Renderer) RenderModel(model interface{}, err error) error {
	if err != nil && !errors.Is(err, models.UserSpaceError{}) {
		return err
	}

	var renderErr error

	switch r.outputType {
	case models.OutputTypeJSON:
		renderErr = r.renderJSON(model)
	case models.OutputTypeASCII:
		renderErr = r.renderASCII(model)
	default:
		panic(fmt.Sprintf("failed to render: unknown output type: %s", r.outputType))
	}

	if renderErr != nil {
		return fmt.Errorf("failed to render model: %w", renderErr)
	}

	return err
}

func (r *Renderer) renderASCII(model interface{}) error {
	templateName := fmt.Sprintf("%T", model)
	templateBuffer, exist := r.asciiTemplates[templateName]

	if !exist {
		return fmt.Errorf("ascii template for model '%s' not exist", templateName)
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
			fnDir:        r.asciiPathDirectory,
			fnPlus:       r.asciiPlus,
			fnMinus:      r.asciiMinus,
		}).
		Parse(
			preprocessRawASCIITemplate(templateBuffer),
		)
	if err != nil {
		return fmt.Errorf("failed to render ascii view '%s': %w", templateName, err)
	}

	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, model)
	if err != nil {
		return fmt.Errorf("failed to execute template '%s': %w", templateName, err)
	}

	fmt.Println(buffer.String())
	return nil
}

func (r *Renderer) renderJSON(model interface{}) error {
	var jsonBuffer []byte
	var marshalErr error

	wrapperModel := struct {
		Type    string      `json:"Type"`
		Payload interface{} `json:"Payload"`
	}{
		Type:    fmt.Sprintf("%T", model),
		Payload: model,
	}

	if r.outputJSONOneLine {
		jsonBuffer, marshalErr = json.Marshal(wrapperModel)
	} else {
		jsonBuffer, marshalErr = json.MarshalIndent(wrapperModel, "", "  ")
	}

	if marshalErr != nil {
		return fmt.Errorf("failed to marshal payload '%v' to json: %w", model, marshalErr)
	}

	fmt.Println(string(jsonBuffer))
	return nil
}
