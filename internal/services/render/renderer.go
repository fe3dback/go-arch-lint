package render

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
	fnLinePrefix = "linePrefix"
	fnDir        = "dir"
	fnPlus       = "plus"
	fnMinus      = "minus"
	fnConcat     = "concat"
)

type (
	Renderer struct {
		colorPrinter      colorPrinter
		referenceRender   referenceRender
		outputType        models.OutputType
		outputJSONOneLine bool
		asciiTemplates    map[string]string
	}
)

func NewRenderer(
	colorPrinter colorPrinter,
	referenceRender referenceRender,
	outputType models.OutputType,
	outputJSONOneLine bool,
	asciiTemplates map[string]string,
) *Renderer {
	return &Renderer{
		colorPrinter:      colorPrinter,
		referenceRender:   referenceRender,
		outputType:        outputType,
		outputJSONOneLine: outputJSONOneLine,
		asciiTemplates:    asciiTemplates,
	}
}

func (r *Renderer) RenderModel(model interface{}, err error) error {
	if err != nil && !errors.Is(err, models.UserSpaceError{}) {
		var referableErr models.ReferableErr
		if errors.As(err, &referableErr) {
			codePreview := r.referenceRender.SourceCode(referableErr.Reference().ExtendRange(1, 1), true, true)
			fmt.Printf("ERR: %s\n", err.Error())
			fmt.Printf("------------\n")
			fmt.Printf("%s\n", codePreview)
		}

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

	modelType, err := r.extractModelType(model)
	if err != nil {
		return fmt.Errorf("failed extract model type from '%T' (maybe not matched pattern: 'CmdXXXOut') : %w", model, err)
	}

	wrapperModel := struct {
		Type    string      `json:"Type"`
		Payload interface{} `json:"Payload"`
	}{
		Type:    modelType,
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

// Rename "anypackage.CmdXXXXOut" to "models.XXXX"
// for back compatible with previous response version
func (r *Renderer) extractModelType(model any) (string, error) {
	const expectedPrefix = "Cmd"
	const expectedSuffix = "Out"

	alias := fmt.Sprintf("%T", model)
	dotIndex := strings.Index(alias, ".")

	if dotIndex == -1 {
		return "", fmt.Errorf("DTO type '%s' without package name", alias)
	}

	dtoName := alias[dotIndex+1:]

	if !strings.HasPrefix(dtoName, expectedPrefix) {
		return "", fmt.Errorf("DTO name '%s' alias '%s' should has prefix '%s'", dtoName, alias, expectedPrefix)
	}

	if !strings.HasSuffix(dtoName, expectedSuffix) {
		return "", fmt.Errorf("DTO name '%s' alias '%s' should has suffix '%s'", dtoName, alias, expectedSuffix)
	}

	return fmt.Sprintf(
		"models.%s",
		strings.TrimPrefix(strings.TrimSuffix(dtoName, expectedSuffix), expectedPrefix),
	), nil
}
