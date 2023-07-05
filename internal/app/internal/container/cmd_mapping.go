package container

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/commands/mapping"
	"github.com/fe3dback/go-arch-lint/internal/models"
)

func (c *Container) ProvideMappingCommand() *cobra.Command {
	return c.provideMappingCommandAssembler().Assemble()
}

func (c *Container) provideMappingCommandAssembler() *mapping.CommandAssembler {
	return mapping.NewMappingCommandAssembler(
		c.provideProjectInfoAssembler(),
		func(ctx context.Context, input models.FlagsMapping) error {
			return c.ProvideRenderer().RenderModel(
				c.provideOperationMapping(input).Behave(ctx, input.Scheme),
			)
		},
	)
}
