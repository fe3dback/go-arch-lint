package container

import (
	"github.com/fe3dback/go-arch-lint/v4/internal/services/colorizer"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer/ascii"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer/json"
	"github.com/fe3dback/go-arch-lint/v4/internal/view"
)

func (c *Container) serviceRenderer() *renderer.Renderer {
	return once(func() *renderer.Renderer {
		return renderer.New(
			json.NewRenderer(),
			ascii.NewRenderer(
				c.serviceAsciiColorizer(),
				view.Templates,
			),
		)
	})
}

func (c *Container) serviceAsciiColorizer() *colorizer.ASCII {
	return once(func() *colorizer.ASCII {
		return colorizer.New()
	})
}
