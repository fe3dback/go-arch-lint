package schema

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"

	"github.com/spf13/cobra"
)

func (c *CommandAssembler) assembleFlags(cmd *cobra.Command) {
	c.withVersion(cmd)
}

func (c *CommandAssembler) withVersion(cmd *cobra.Command) {
	cmd.PersistentFlags().IntP(
		flagVersion,
		"",
		c.localFlags.Version,
		fmt.Sprintf("json schema version to output (min: %d, max: %d)",
			1,
			models.SupportedVersion,
		),
	)
}
