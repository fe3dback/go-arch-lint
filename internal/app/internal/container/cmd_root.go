package container

import (
	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/commands/root"
	"github.com/fe3dback/go-arch-lint/internal/models"
)

func (c *Container) ProvideRootCommand() *cobra.Command {
	assembler := c.provideRootCommandAssembler()
	return assembler.Assemble()
}

func (c *Container) provideRootCommandAssembler() *root.CommandAssembler {
	return root.NewRootCommandAssembler(
		c.provideRootCommandFlagsAssemblingFn(),
		c.provideCommands(),
	)
}

func (c *Container) provideCommands() []*cobra.Command {
	return []*cobra.Command{
		c.ProvideVersionCommand(),
		c.ProvideCheckCommand(),
	}
}

func (c *Container) provideRootCommandFlagsAssemblingFn() func(flags models.FlagsRoot) error {
	return func(flags models.FlagsRoot) error {
		c.flags = flags
		return nil
	}
}
