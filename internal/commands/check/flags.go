package check

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"

	"github.com/spf13/cobra"
)

func (c *CommandAssembler) assembleFlags(cmd *cobra.Command) {
	c.withProjectPath(cmd)
	c.withArchFileName(cmd)
	c.withMaxWarnings(cmd)
}

func (c *CommandAssembler) withProjectPath(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		flagProjectPath,
		"",
		c.localFlags.ProjectPath,
		fmt.Sprintf("absolute path to project directory (where '%s' is located)", models.ProjectInfoDefaultArchFileName),
	)
}

func (c *CommandAssembler) withArchFileName(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		flagArchFile,
		"",
		c.localFlags.ArchFile,
		"arch file path",
	)
}

func (c *CommandAssembler) withMaxWarnings(cmd *cobra.Command) {
	cmd.PersistentFlags().IntP(
		flagMaxWarnings,
		"",
		c.localFlags.MaxWarnings,
		"max number of warnings to output",
	)
}
