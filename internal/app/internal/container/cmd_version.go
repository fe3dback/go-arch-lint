package container

import (
	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/commands/version"
)

func (c *Container) ProvideVersionCommand() *cobra.Command {
	return c.provideVersionCommandAssembler().Assemble()
}

func (c *Container) provideVersionCommandAssembler() *version.CommandAssembler {
	return version.NewVersionCommandAssembler(func() error {
		return c.ProvideRenderer().RenderModel(
			c.provideOperationVersion().Behave(),
		)
	})
}
