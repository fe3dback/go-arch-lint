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
	return func(context *cli.Context) error {
		model, err := handler.Execute(context)
		if err != nil {
			return fmt.Errorf("command '%s' failed: %w", name, err)
		}

		rnd := c.serviceRenderer()
		// todo: render mode
		out, err := rnd.Render(models.OutputTypeDefault, model)
		if err != nil {
			return fmt.Errorf("command '%s' render failed: %w", name, err)
		}

		fmt.Println(out)
		return nil
	}
}
