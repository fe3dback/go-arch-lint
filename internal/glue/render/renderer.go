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
	fnColorize = "colorize"
)

type (
	Renderer struct {
		colorPrinter      ColorPrinter
		outputType        models.OutputType
		outputJsonOneLine bool
		asciiTemplates    map[string]string
	}
)

func NewRenderer(
	colorPrinter ColorPrinter,
	outputType models.OutputType,
	outputJsonOneLine bool,
	asciiTemplates map[string]string,
) *Renderer {
	return &Renderer{
		colorPrinter:      colorPrinter,
		outputType:        outputType,
		outputJsonOneLine: outputJsonOneLine,
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
		renderErr = r.renderJson(model)
	case models.OutputTypeASCII:
		renderErr = r.renderAscii(model)
	default:
		panic(fmt.Sprintf("failed to render: unknown output type: %s", r.outputType))
	}

	if renderErr != nil {
		return fmt.Errorf("failed to render model: %w", renderErr)
	}

	return err
}

func (r *Renderer) renderAscii(model interface{}) error {
	templateName := fmt.Sprintf("%T", model)
	templateBuffer, exist := r.asciiTemplates[templateName]

	if !exist {
		return fmt.Errorf("ascii template for model '%s' not exist", templateName)
	}

	tpl, err := template.
		New(templateName).
		Funcs(map[string]interface{}{
			fnColorize: r.asciiColorize,
		}).
		Parse(
			preprocessRawAsciiTemplate(templateBuffer),
		)
	if err != nil {
		return fmt.Errorf("failed to render ascii view '%s': %s", templateName, err)
	}

	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, model)
	if err != nil {
		return fmt.Errorf("failed to execute template '%s': %s", templateName, err)
	}

	fmt.Println(buffer.String())
	return nil
}

func (r *Renderer) asciiColorize(color string, value interface{}) (string, error) {
	colorizer := newColorizer(r.colorPrinter)
	out, err := colorizer.colorize(
		color,
		fmt.Sprintf("%s", value),
	)
	if err != nil {
		return "", fmt.Errorf("failed colorize: %s", err)
	}

	return out, nil
}

func (r *Renderer) renderJson(model interface{}) error {
	var jsonBuffer []byte
	var marshalErr error

	wrapperModel := struct {
		Type    string      `json:"Type"`
		Payload interface{} `json:"Payload"`
	}{
		Type:    fmt.Sprintf("%T", model),
		Payload: model,
	}

	if r.outputJsonOneLine {
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
