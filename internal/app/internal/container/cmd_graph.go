package container

import (
	"context"

	"github.com/fe3dback/go-arch-lint/internal/commands/graph"
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/spf13/cobra"
)

func (c *Container) ProvideGraphCommand() *cobra.Command {
	return c.provideGraphCommandAssembler().Assemble()
}

func (c *Container) provideGraphCommandAssembler() *graph.CommandAssembler {
	return graph.NewGraphCommandAssembler(
		c.provideProjectInfoAssembler(),
		func(ctx context.Context, input models.FlagsGraph) error {
			return c.ProvideRenderer().RenderModel(
				c.provideOperationGraph(input).Behave(ctx, input),
			)
		},
	)
}
