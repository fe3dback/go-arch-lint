package renderer

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Renderer struct {
	jsonRenderer  typeRenderer
	asciiRenderer typeRenderer
}

func New(
	jsonRenderer typeRenderer,
	asciiRenderer typeRenderer,
) *Renderer {
	return &Renderer{
		jsonRenderer:  jsonRenderer,
		asciiRenderer: asciiRenderer,
	}
}

func (r *Renderer) Render(outputType models.OutputType, model any) (string, error) {
	if outputType == models.OutputTypeDefault {
		outputType = models.OutputTypeASCII
	}

	switch outputType {
	case models.OutputTypeJSON:
		return r.jsonRenderer.Render(model)
	case models.OutputTypeASCII:
		return r.asciiRenderer.Render(model)
	default:
		return "", fmt.Errorf("unknown renderer type: %s", outputType)
	}
}
