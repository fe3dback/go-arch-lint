package container

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type CommandHandler interface {
	Execute(ctx *cli.Context) (any, error)
}

func (c *Container) makeBeforeCode() cli.BeforeFunc {
	return func(context *cli.Context) error {
		c.cCtx = context
		return nil
	}
}

func (c *Container) makeCliCommand(name string, factory func() CommandHandler) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		handler := factory()
		model, err := handler.Execute(cCtx)
		if err != nil {
			renderErr := c.render(cCtx, name, true, c.serviceErrorBuilder().BuildError(err))
			if renderErr != nil {
				return renderErr
			}

			return models.NewUserLandError(err)
		}

		return c.render(cCtx, name, false, model)
	}
}

func (c *Container) render(cCtx *cli.Context, name string, isErr bool, model any) error {
	outputType, err := extractOutputType(cCtx)
	if err != nil {
		return fmt.Errorf("failed extracting output type: %w", err)
	}

	renderMode := models.RenderOptions{
		OutputType: outputType,
		FormatJson: !cCtx.Bool(models.FlagOutputJSONWithoutFormatting),
	}

	rnd := c.serviceRenderer()
	out, err := rnd.Render(model, renderMode)
	if err != nil {
		return fmt.Errorf("command '%s' render failed: %w", name, err)
	}

	if isErr {
		_, err = fmt.Fprintln(cCtx.App.Writer, out)
	} else {
		_, err = fmt.Fprintln(cCtx.App.ErrWriter, out)
	}

	if err != nil {
		return fmt.Errorf("print to stdout failed: %w", err)
	}

	return nil
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