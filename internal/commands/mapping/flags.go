package mapping

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

func (c *CommandAssembler) assembleFlags(cmd *cobra.Command) {
	c.withProjectPath(cmd)
	c.withArchFileName(cmd)
	c.withScheme(cmd)
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

func (c *CommandAssembler) withScheme(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		flagScheme,
		"s",
		c.localFlags.Scheme,
		fmt.Sprintf(
			"display scheme [%s]",
			strings.Join(validSchemes, ","),
		),
	)
}
