package root

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

func (c *CommandAssembler) assembleFlags(cmd *cobra.Command) {
	c.withColors(cmd)
	c.withOutputType(cmd)
	c.withOutputJsonOneLine(cmd)
	c.withJsonAlias(cmd)
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

func (c *CommandAssembler) withOutputJsonOneLine(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(
		flagOutputJsonOneLine,
		"",
		false,
		"format JSON as single line payload (without line breaks), only for json output type",
	)
}

func (c *CommandAssembler) withJsonAlias(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(
		flagAliasJson,
		"",
		false,
		fmt.Sprintf("(alias for --%s=%s)",
			flagOutputType,
			models.OutputTypeJSON,
		),
	)
}
