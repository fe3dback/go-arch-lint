package container

import (
	"github.com/urfave/cli/v2"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/colorizer"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/config/reader"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/config/reader/yaml"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer/ascii"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer/json"
	"github.com/fe3dback/go-arch-lint/v4/internal/view"
)

func (c *Container) serviceRenderer(cCtx *cli.Context) *renderer.Renderer {
	return once(func() *renderer.Renderer {
		return renderer.New(
			json.NewRenderer(),
			ascii.NewRenderer(
				c.serviceAsciiColorizer(cCtx),
				view.Templates,
			),
		)
	})
}

func (c *Container) serviceAsciiColorizer(cCtx *cli.Context) *colorizer.ASCII {
	return once(func() *colorizer.ASCII {
		return colorizer.New(cCtx.Bool(models.FlagOutputUseAsciiColors))
	})
}

func (c *Container) serviceConfigReader() *reader.Reader {
	return once(func() *reader.Reader {
		return reader.NewReader(
			c.serviceConfigReaderYAML(),
		)
	})
}

func (c *Container) serviceConfigReaderYAML() *yaml.Reader {
	return once(func() *yaml.Reader {
		return yaml.NewReader()
	})
}
