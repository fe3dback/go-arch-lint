package container

import (
	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/commands/selfInspect"
	"github.com/fe3dback/go-arch-lint/internal/models"
)

func (c *Container) ProvideSelfInspectCommand() *cobra.Command {
	return c.provideSelfInspectCommandAssembler().Assemble()
}

func (c *Container) provideSelfInspectCommandAssembler() *selfInspect.CommandAssembler {
	return selfInspect.NewSelfInspectCommandAssembler(
		c.provideProjectInfoAssembler(),
		func(input models.FlagsSelfInspect) error {
			return c.ProvideRenderer().RenderModel(c.provideOperationSelfInspect(input).Behave())
		})
}
