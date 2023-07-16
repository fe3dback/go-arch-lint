package container

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/operations/schema"
	"github.com/spf13/cobra"
)

func (c *Container) commandSchema() (*cobra.Command, runner) {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "json schema for arch file inspection",
		Long:  "useful for integrations with ide's and editor plugins",
	}

	in := models.CmdSchemaIn{
		Version: 0,
	}

	cmd.PersistentFlags().IntVar(&in.Version, "version", in.Version, fmt.Sprintf("json schema version to output (min: %d, max: %d)",
		models.SupportedVersionMin,
		models.SupportedVersionMax,
	))

	return cmd, func(act *cobra.Command) (any, error) {
		return c.commandSchemaOperation().Behave(in)
	}
}

func (c *Container) commandSchemaOperation() *schema.Operation {
	return schema.NewOperation(
		c.provideJsonSchemaProvider(),
	)
}
