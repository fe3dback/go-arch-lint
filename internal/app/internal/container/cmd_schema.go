package container

import (
	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/commands/schema"
	"github.com/fe3dback/go-arch-lint/internal/models"
)

func (c *Container) ProvideSchemaCommand() *cobra.Command {
	assembler := c.provideSchemaCommandAssembler()
	return assembler.Assemble()
}

func (c *Container) provideSchemaCommandAssembler() *schema.CommandAssembler {
	return schema.NewSchemaCommandAssembler(
		c.provideJsonSchemaProvider(),
		func(input models.FlagsSchema) error {
			return c.ProvideRenderer().RenderModel(
				c.provideSchemaService().Behave(input),
			)
		},
	)
}
