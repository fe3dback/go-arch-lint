package renderer

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Renderer struct {
	jsonRenderer  jsonRenderer
	asciiRenderer asciiRenderer
}

func New(
	jsonRenderer jsonRenderer,
	asciiRenderer asciiRenderer,
) *Renderer {
	return &Renderer{
		jsonRenderer:  jsonRenderer,
		asciiRenderer: asciiRenderer,
	}
}

func (r *Renderer) Render(model any, options models.RenderOptions) (string, error) {
	switch options.OutputType {
	case models.OutputTypeJSON:
		return r.jsonRenderer.Render(model, options.FormatJson)
	case models.OutputTypeASCII:
		return r.asciiRenderer.Render(model)
	default:
		return "", fmt.Errorf("unknown renderer type: %s", options.OutputType)
	}
}
