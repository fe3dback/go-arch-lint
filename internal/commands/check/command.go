package check

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

const (
	flagMaxWarnings = "max-warnings"
	flagProjectPath = "project-path"
	flagArchFile    = "arch-file"
)

const defaultArchFileName = ".go-arch-lint.yml"

type (
	processorFn = func() error

	CommandAssembler struct {
		processorFn processorFn
		flags       *models.FlagsCheck
	}
)

func NewCheckCommandAssembler(processorFn processorFn) *CommandAssembler {
	return &CommandAssembler{
		processorFn: processorFn,
		flags: &models.FlagsCheck{
			MaxWarnings: 512,
			ProjectPath: "",
			ArchFile:    defaultArchFileName,
		},
	}
}

func (c *CommandAssembler) Assemble() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "check",
		Short:             "check project architecture by yaml file",
		Long:              "compare project *.go files with arch defined in spec file",
		PersistentPreRunE: c.prePersist,
		RunE:              c.invoke,
	}

	c.assembleFlags(cmd)

	return cmd
}

func (c *CommandAssembler) invoke(_ *cobra.Command, _ []string) error {
	fmt.Printf("%+v", c.flags)

	return c.processorFn()
}

func (c *CommandAssembler) prePersist(cmd *cobra.Command, _ []string) error {
	rootDirectory, err := cmd.Flags().GetString(flagProjectPath)
	if err != nil {
		return c.failedToGetFlag(err, flagProjectPath)
	}

	archFile, err := cmd.Flags().GetString(flagArchFile)
	if err != nil {
		return c.failedToGetFlag(err, flagArchFile)
	}

	maxWarnings, err := cmd.Flags().GetInt(flagMaxWarnings)
	if err != nil {
		return c.failedToGetFlag(err, flagMaxWarnings)
	}

	const warningsRangeMin = 1
	const warningsRangeMax = 32768

	if maxWarnings < warningsRangeMin || maxWarnings > warningsRangeMax {
		return fmt.Errorf(
			"flag '%s' should by in range [%d .. %d]",
			flagMaxWarnings,
			warningsRangeMin,
			warningsRangeMax,
		)
	}

	// assemble flags
	c.flags.ProjectPath = rootDirectory
	c.flags.ArchFile = archFile
	c.flags.MaxWarnings = maxWarnings

	return nil
}
