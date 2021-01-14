package root

import (
	"fmt"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"

	"github.com/spf13/cobra"
)

func (c *CommandAssembler) assembleFlags(cmd *cobra.Command) {
	c.withColors(cmd)
	c.withOutputType(cmd)
	c.withOutputJSONOneLine(cmd)
	c.withJSONAlias(cmd)
}

func (c *CommandAssembler) withColors(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagUseColors, "", true, "use ANSI colors in terminal output")
}

func (c *CommandAssembler) withOutputType(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		flagOutputType,
		"",
		models.OutputTypeDefault,
		fmt.Sprintf(
			"type of command output, variants: [%s]",
			strings.Join(models.OutputTypeVariantsConst, ", "),
		),
	)
}

func (c *CommandAssembler) withOutputJSONOneLine(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(
		flagOutputJSONOneLine,
		"",
		false,
		"format JSON as single line payload (without line breaks), only for json output type",
	)
}

func (c *CommandAssembler) withJSONAlias(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(
		flagAliasJSON,
		"",
		false,
		fmt.Sprintf("(alias for --%s=%s)",
			flagOutputType,
			models.OutputTypeJSON,
		),
	)
}
