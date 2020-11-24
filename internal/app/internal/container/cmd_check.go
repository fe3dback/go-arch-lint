package container

import (
	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/commands/check"
)

func (c *Container) ProvideCheckCommand() *cobra.Command {
	assembler := c.provideCheckCommandAssembler()
	return assembler.Assemble()
}

func (c *Container) provideCheckCommandAssembler() *check.CommandAssembler {
	return check.NewCheckCommandAssembler(func() error {
		return nil

		//renderer := c.ProvideRenderer()
		//return renderer.RenderModel(
		//	c.provideVersionService().Behave(),
		//)
	})
}
