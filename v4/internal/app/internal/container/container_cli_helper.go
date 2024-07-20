package container

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type CommandHandler interface {
	Execute(ctx *cli.Context) (any, error)
}

func (c *Container) makeCliCommand(name string, handler CommandHandler) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		model, err := handler.Execute(cCtx)
		if err != nil {
			return fmt.Errorf("command '%s' failed: %w", name, err)
		}

		outputType, err := extractOutputType(cCtx)
		if err != nil {
			return fmt.Errorf("failed extracting output type: %w", err)
		}

		rnd := c.serviceRenderer(cCtx)
		out, err := rnd.Render(model, models.RenderOptions{
			OutputType: outputType,
			FormatJson: !cCtx.Bool(models.FlagOutputJSONWithoutFormatting),
		})
		if err != nil {
			return fmt.Errorf("command '%s' render failed: %w", name, err)
		}

		_, err = fmt.Fprintln(cCtx.App.Writer, out)
		if err != nil {
			return err
		}

		return nil
	}
}

func extractOutputType(cCtx *cli.Context) (models.OutputType, error) {
	forceJSON := cCtx.Bool(models.FlagOutputTypeJSON)
	outputType := cCtx.String(models.FlagOutputType)

	if outputType == models.OutputTypeDefault {
		outputType = models.OutputTypeASCII
	}

	if cCtx.IsSet(models.FlagOutputTypeJSON) && cCtx.IsSet(models.FlagOutputType) {
		if forceJSON && outputType == models.OutputTypeASCII {
			return "", fmt.Errorf("flag '--%s' not compatible with '--%s %s'",
				models.FlagOutputTypeJSON,
				models.FlagOutputType,
				models.OutputTypeASCII,
			)
		}
	}

	if forceJSON {
		outputType = models.OutputTypeJSON
	}

	if cCtx.IsSet(models.FlagOutputJSONWithoutFormatting) && outputType != models.OutputTypeJSON {
		return "", fmt.Errorf("flag '--%s' used only with '--%s'",
			models.FlagOutputJSONWithoutFormatting,
			models.FlagOutputTypeJSON,
		)
	}

	if cCtx.IsSet(models.FlagOutputUseAsciiColors) && outputType != models.OutputTypeASCII {
		return "", fmt.Errorf("flag '--%s' used only with '--%s %s'",
			models.FlagOutputUseAsciiColors,
			models.FlagOutputType,
			models.OutputTypeASCII,
		)
	}

	return outputType, nil
}
