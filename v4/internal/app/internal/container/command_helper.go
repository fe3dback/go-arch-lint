package container

import (
	"errors"
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
			renderErr := c.render(cCtx, name, c.extractStdoutError(cCtx, name, err))
			if renderErr != nil {
				return renderErr
			}

			return models.NewUserLandError(err)
		}

		return c.render(cCtx, name, model)
	}
}

func (c *Container) extractStdoutError(cCtx *cli.Context, name string, err error) models.CmdStdoutErrorOut {
	ref := models.NewInvalidReference()

	refError := models.ReferencedError{}
	if errors.As(err, &refError) {
		ref = refError.Reference()
	}

	preview := ""
	if ref.Valid {
		preview, _ = c.servicePrinter().Print(ref, models.CodePrintOpts{
			LineNumbers: true,
			Arrows:      true,
			Highlight:   cCtx.Bool(models.FlagOutputUseAsciiColors),
			Mode:        models.CodePrintModeExtend,
		})
	}

	return models.CmdStdoutErrorOut{
		Error:            fmt.Errorf("command '%s' failed: %w", name, err).Error(),
		Reference:        ref,
		ReferencePreview: preview,
	}
}

func (c *Container) render(cCtx *cli.Context, name string, model any) error {
	outputType, err := extractOutputType(cCtx)
	if err != nil {
		return fmt.Errorf("failed extracting output type: %w", err)
	}

	renderMode := models.RenderOptions{
		OutputType: outputType,
		FormatJson: !cCtx.Bool(models.FlagOutputJSONWithoutFormatting),
	}

	rnd := c.serviceRenderer(cCtx)
	out, err := rnd.Render(model, renderMode)
	if err != nil {
		return fmt.Errorf("command '%s' render failed: %w", name, err)
	}

	_, err = fmt.Fprintln(cCtx.App.Writer, out)
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
