package root

import (
	"github.com/fe3dback/go-arch-lint/internal/models"

	"github.com/spf13/cobra"
)

const (
	flagUseColors         = "output-color"
	flagOutputType        = "output-type"
	flagOutputJSONOneLine = "output-json-one-line"
	flagAliasJSON         = "json"
)

type (
	flagAssemblingFn = func(root models.FlagsRoot) error

	CommandAssembler struct {
		flagsAssemblingFn flagAssemblingFn
		commands          []*cobra.Command
		flags             *models.FlagsRoot
	}
)

func NewRootCommandAssembler(
	flagsAssemblingFn flagAssemblingFn,
	commands []*cobra.Command,
) *CommandAssembler {
	return &CommandAssembler{
		flagsAssemblingFn: flagsAssemblingFn,
		commands:          commands,
		flags: &models.FlagsRoot{
			UseColors:         true,
			OutputType:        models.OutputTypeDefault,
			OutputJsonOneLine: false,
		},
	}
}

func (c *CommandAssembler) Assemble() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "go-arch-lint",
		Short:             "Golang architecture linter",
		Long:              "Check all project imports and compare to arch rules defined in yaml file.\nRead full documentation in: https://github.com/fe3dback/go-arch-lint",
		PersistentPreRunE: c.prePersist,
		RunE:              c.invoke,
		SilenceErrors:     true, // redirect to stderr
		SilenceUsage:      true,
	}

	// apply root args
	c.assembleFlags(cmd)

	// apply sub commands
	for _, command := range c.commands {
		if err := assertCommandIsValid(command); err != nil {
			panic(err)
		}

		cmd.AddCommand(command)
	}

	return cmd
}

func (c *CommandAssembler) invoke(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func (c *CommandAssembler) prePersist(cmd *cobra.Command, _ []string) error {
	useColors, err := cmd.Flags().GetBool(flagUseColors)
	if err != nil {
		return failedToGetFlag(err, flagUseColors)
	}

	outputType, err := c.prepareOutputType(cmd)
	if err != nil {
		return err
	}

	outputJsonOneLine, err := cmd.Flags().GetBool(flagOutputJSONOneLine)
	if err != nil {
		return failedToGetFlag(err, flagOutputJSONOneLine)
	}

	// all root cmd flags is global, and we should set it to global container context
	c.flags.UseColors = useColors
	c.flags.OutputType = outputType
	c.flags.OutputJsonOneLine = outputJsonOneLine

	return c.flagsAssemblingFn(*c.flags)
}
