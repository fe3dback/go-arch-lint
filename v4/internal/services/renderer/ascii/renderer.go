package ascii

import (
	"fmt"
)

type ASCII struct {
	renderer renderer
}

func NewRenderer(
	renderer renderer,
	templates map[string][]byte,
) (*ASCII, error) {
	for id, templateBody := range templates {
		err := renderer.RegisterTemplate(id, templateBody)
		if err != nil {
			return nil, fmt.Errorf("failed create renderer: %w", err)
		}
	}

	return &ASCII{
		renderer: renderer,
	}, nil
}

func (r *ASCII) Render(model any) (string, error) {
	return r.renderer.Render(
		fmt.Sprintf("%T", model),
		model,
	)
}
