package container

import (
	"github.com/fe3dback/go-arch-lint-sdk/pkg/tpl"
	"github.com/fe3dback/go-arch-lint-sdk/pkg/tpl/colorizer"
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

func (c *Container) sdkRenderer() *tpl.Renderer {
	return once(func() *tpl.Renderer {
		return tpl.NewRenderer(
			c.sdkColorizer(),
		)
	})
}

func (c *Container) sdkColorizer() *colorizer.Colorizer {
	return once(func() *colorizer.Colorizer {
		return colorizer.New(c.cCtx.Bool(models.FlagOutputUseAsciiColors), false)
	})
}
