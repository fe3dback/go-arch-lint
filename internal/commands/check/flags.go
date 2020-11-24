package check

import (
	"fmt"

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
		c.flags.ProjectPath,
		fmt.Sprintf("absolute path to project directory (where '%s' is located)", defaultArchFileName),
	)
}

func (c *CommandAssembler) withArchFileName(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		flagArchFile,
		"",
		c.flags.ArchFile,
		"arch file path",
	)
}

func (c *CommandAssembler) withMaxWarnings(cmd *cobra.Command) {
	cmd.PersistentFlags().IntP(
		flagMaxWarnings,
		"",
		c.flags.MaxWarnings,
		"max number of warnings to output",
	)
}
