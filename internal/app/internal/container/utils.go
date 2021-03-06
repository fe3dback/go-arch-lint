package container

import (
	"github.com/logrusorgru/aurora/v3"

	"github.com/fe3dback/go-arch-lint/internal/glue/printer"
	"github.com/fe3dback/go-arch-lint/internal/glue/render"
	"github.com/fe3dback/go-arch-lint/internal/view"
)

func (c *Container) provideColorPrinter() *printer.ColorPrinter {
	return printer.NewColorPrinter(
		c.provideAurora(),
	)
}

func (c *Container) provideAurora() aurora.Aurora {
	return aurora.NewAurora(
		c.flags.UseColors,
	)
}

func (c *Container) ProvideRenderer() *render.Renderer {
	return render.NewRenderer(
		c.provideColorPrinter(),
		c.flags.OutputType,
		c.flags.OutputJsonOneLine,
		view.Templates,
	)
}
