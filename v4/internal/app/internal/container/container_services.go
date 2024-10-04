package container

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/pkg/tpl/codeprinter"
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer/ascii"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/renderer/json"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/spec"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/xstdout"
	"github.com/fe3dback/go-arch-lint/v4/internal/view"
)

func (c *Container) serviceSpecFetcher() *spec.Fetcher {
	return once(func() *spec.Fetcher {
		return spec.NewFetcher(
			c.sdk(),
			arch.PathRelative(c.cCtx.String(models.FlagArchConfigRelativePath)),
		)
	})
}

func (c *Container) serviceErrorBuilder() *xstdout.ErrorBuilder {
	return once(func() *xstdout.ErrorBuilder {
		return xstdout.NewErrorBuilder(
			c.servicePrinter(),
		)
	})
}

func (c *Container) serviceRenderer() *renderer.Renderer {
	return once(func() *renderer.Renderer {
		asciiRenderer, err := ascii.NewRenderer(
			c.sdkRenderer(),
			view.Templates,
		)
		if err != nil {
			panic(fmt.Errorf("failed create ascii renderer: %w", err))
		}

		return renderer.New(
			json.NewRenderer(),
			asciiRenderer,
		)
	})
}

func (c *Container) servicePrinter() *codeprinter.Printer {
	return once(func() *codeprinter.Printer {
		return codeprinter.NewPrinter(
			c.servicePrinterExtractorRaw(),
			c.servicePrinterExtractorHL(),
			c.colorEnv() == arch.TerminalColorEnvColored,
		)
	})
}

func (c *Container) servicePrinterExtractorRaw() *codeprinter.ExtractorRaw {
	return once(codeprinter.NewExtractorRaw)
}

func (c *Container) servicePrinterExtractorHL() *codeprinter.ExtractorHL {
	return once(codeprinter.NewExtractorHL)
}
