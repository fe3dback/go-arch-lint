package check

import (
	"context"
	"fmt"

	terminal "github.com/fe3dback/span-terminal"

	"github.com/fe3dback/go-arch-lint/internal/models"

	"github.com/spf13/cobra"
)

const (
	flagMaxWarnings = "max-warnings"
	flagProjectPath = "project-path"
	flagArchFile    = "arch-file"
)

type (
	processorFn = func(context.Context, models.FlagsCheck) error

	CommandAssembler struct {
		projectInfoAssembler ProjectInfoAssembler
		processorFn          processorFn
		localFlags           *localFlags
	}

	localFlags struct {
		MaxWarnings int
		ProjectPath string
		ArchFile    string
	}
)

func NewCheckCommandAssembler(
	projectInfoAssembler ProjectInfoAssembler,
	processorFn processorFn,
) *CommandAssembler {
	return &CommandAssembler{
		projectInfoAssembler: projectInfoAssembler,
		processorFn:          processorFn,
		localFlags: &localFlags{
			MaxWarnings: 512,
			ProjectPath: "./",
			ArchFile:    models.ProjectInfoDefaultArchFileName,
		},
	}
}

func (c *CommandAssembler) Assemble() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "check",
		Aliases: []string{"c"},
		Short:   "check project architecture by yaml file",
		Long:    "compare project *.go files with arch defined in spec file",
		PreRunE: c.prePersist,
		RunE:    c.invoke,
	}

	c.assembleFlags(cmd)

	return cmd
}

func (c *CommandAssembler) invoke(cmd *cobra.Command, _ []string) error {
	// track progress of this command, with stdout drawing
	terminal.CaptureOutput()
	defer terminal.ReleaseOutput()

	// run
	input, err := c.assembleInput()
	if err != nil {
		return fmt.Errorf("failed to assemble input params: %w", err)
	}

	return c.processorFn(cmd.Context(), input)
}

func (c *CommandAssembler) prePersist(cmd *cobra.Command, _ []string) error {
	rootDirectory, err := cmd.Flags().GetString(flagProjectPath)
	if err != nil {
		return failedToGetFlag(err, flagProjectPath)
	}

	archFile, err := cmd.Flags().GetString(flagArchFile)
	if err != nil {
		return failedToGetFlag(err, flagArchFile)
	}

	maxWarnings, err := cmd.Flags().GetInt(flagMaxWarnings)
	if err != nil {
		return failedToGetFlag(err, flagMaxWarnings)
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

	// assemble localFlags
	c.localFlags.ProjectPath = rootDirectory
	c.localFlags.ArchFile = archFile
	c.localFlags.MaxWarnings = maxWarnings

	return nil
}

func (c *CommandAssembler) assembleInput() (models.FlagsCheck, error) {
	projectInfo, err := c.projectInfoAssembler.ProjectInfo(
		c.localFlags.ProjectPath,
		c.localFlags.ArchFile,
	)
	if err != nil {
		return models.FlagsCheck{}, fmt.Errorf("failed to assemble project info: %w", err)
	}

	return models.FlagsCheck{
		Project:     projectInfo,
		MaxWarnings: c.localFlags.MaxWarnings,
	}, nil
}
