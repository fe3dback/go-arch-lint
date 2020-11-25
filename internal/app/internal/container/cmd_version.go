package container

import (
	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/commands/version"
)

func (c *Container) ProvideVersionCommand() *cobra.Command {
	assembler := c.provideVersionCommandAssembler()
	return assembler.Assemble()
}

func (c *Container) provideVersionCommandAssembler() *version.CommandAssembler {
	return version.NewVersionCommandAssembler(func() error {
		return c.ProvideRenderer().RenderModel(
			c.provideVersionService().Behave(),
		)
	})
}
