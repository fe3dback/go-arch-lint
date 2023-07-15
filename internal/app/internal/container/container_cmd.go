package container

import (
	"fmt"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/spf13/cobra"
)

type runner = func(cmd *cobra.Command) (any, error)

func (c *Container) CommandRoot() *cobra.Command {
	flags := models.FlagsRoot{
		UseColors:         true,
		OutputType:        models.OutputTypeDefault,
		OutputJsonOneLine: false,
	}
	flagAliasOutputTypeJson := false

	rootCmd := &cobra.Command{
		Use:           "go-arch-lint",
		Short:         "Golang architecture linter",
		Long:          "Check all project imports and compare to arch rules defined in yaml file.\nRead full documentation in: https://github.com/fe3dback/go-arch-lint",
		SilenceErrors: true, // redirect to stderr
		SilenceUsage:  true,
		RunE: func(act *cobra.Command, _ []string) error {
			return act.Help()
		},
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			// alias preprocessor
			if flagAliasOutputTypeJson {
				if flags.OutputType != models.OutputTypeDefault && flags.OutputType != models.OutputTypeJSON {
					return fmt.Errorf("flag --%s not compatible with --%s=%s",
						"json",
						"output-type",
						flags.OutputType,
					)
				}

				flags.OutputType = models.OutputTypeJSON
			}

			// fallback to default's
			if flags.OutputType == models.OutputTypeDefault {
				flags.OutputType = models.OutputTypeASCII
			}

			// validate
			outputTypeIsValid := false
			for _, validValue := range models.OutputTypeValues {
				if flags.OutputType == validValue {
					outputTypeIsValid = true
					break
				}
			}

			if !outputTypeIsValid {
				return fmt.Errorf("unknown output-type: %s", flags.OutputType)
			}

			// save global flags for another child commands
			c.flags = flags
			return nil
		},
	}

	// define global flags
	rootCmd.PersistentFlags().BoolVar(&flags.UseColors, "output-color", flags.UseColors, "use ANSI colors in terminal output")
	rootCmd.PersistentFlags().StringVar(&flags.OutputType, "output-type", flags.OutputType, fmt.Sprintf("type of command output, variants: [%s]", strings.Join(models.OutputTypeValues, ", ")))
	rootCmd.PersistentFlags().BoolVar(&flags.OutputJsonOneLine, "output-json-one-line", flags.OutputJsonOneLine, "format JSON as single line payload (without line breaks), only for json output type")
	rootCmd.PersistentFlags().BoolVar(&flagAliasOutputTypeJson, "json", flagAliasOutputTypeJson, fmt.Sprintf("(alias for --%s=%s)",
		"output-type",
		models.OutputTypeJSON,
	))

	// apply sub commands
	for _, subCmd := range c.commands() {
		if subCmd.PersistentPreRun != nil {
			panic(fmt.Errorf("root sub command '%s' should not have 'PersistentPreRun', "+
				"use 'PreRun' instead", rootCmd.Name(),
			))
		}

		if subCmd.PersistentPreRunE != nil {
			panic(fmt.Errorf("root sub command '%s' should not have 'PersistentPreRunE', "+
				"use 'PreRunE' instead", rootCmd.Name(),
			))
		}

		rootCmd.AddCommand(subCmd)
	}

	return rootCmd
}

func (c *Container) commands() []*cobra.Command {
	type exec struct {
		cmd  *cobra.Command
		runE runner
	}

	unwrap := func(cmd *cobra.Command, r runner) exec {
		return exec{cmd: cmd, runE: r}
	}

	executors := []exec{
		unwrap(c.commandVersion()),
		unwrap(c.commandSelfInspect()),
		unwrap(c.commandSchema()),
		unwrap(c.commandCheck()),
		unwrap(c.commandMapping()),
		unwrap(c.commandGraph()),
	}

	list := make([]*cobra.Command, 0, len(executors))
	for _, x := range executors {
		x := x
		x.cmd.RunE = func(activeCmd *cobra.Command, _ []string) error {
			return c.ProvideRenderer().RenderModel(x.runE(activeCmd))
		}
		list = append(list, x.cmd)
	}

	return list
}
