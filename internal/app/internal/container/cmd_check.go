package container

import (
	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/commands/check"
	"github.com/fe3dback/go-arch-lint/internal/models"
)

func (c *Container) ProvideCheckCommand() *cobra.Command {
	assembler := c.provideCheckCommandAssembler()
	return assembler.Assemble()
}

func (c *Container) provideCheckCommandAssembler() *check.CommandAssembler {
	return check.NewCheckCommandAssembler(func(input models.FlagsCheck) error {
		return c.ProvideRenderer().RenderModel(
			c.provideCheckService(input).Behave(),
		)
	})
}
