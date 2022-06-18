package selfInspect

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

const (
	flagProjectPath = "project-path"
	flagArchFile    = "arch-file"
)

type (
	processorFn = func(input models.FlagsSelfInspect) error

	CommandAssembler struct {
		projectInfoAssembler projectInfoAssembler
		processorFn          processorFn
		localFlags           *localFlags
	}

	localFlags struct {
		ProjectPath string
		ArchFile    string
	}
)

func NewSelfInspectCommandAssembler(projectInfoAssembler projectInfoAssembler, processorFn processorFn) *CommandAssembler {
	return &CommandAssembler{
		projectInfoAssembler: projectInfoAssembler,
		processorFn:          processorFn,
		localFlags: &localFlags{
			ProjectPath: "./",
			ArchFile:    models.ProjectInfoDefaultArchFileName,
		},
	}
}

func (c *CommandAssembler) Assemble() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "self-inspect",
		Short:   "will validate arch config and arch setup",
		Long:    "this useful for IDE plugins and other tool integration",
		PreRunE: c.prePersist,
		RunE:    c.invoke,
	}

	c.assembleFlags(cmd)

	return cmd
}

func (c *CommandAssembler) invoke(_ *cobra.Command, _ []string) error {
	// run
	input, err := c.assembleInput()
	if err != nil {
		return fmt.Errorf("failed to assemble input params: %w", err)
	}

	return c.processorFn(input)
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

	// assemble localFlags
	c.localFlags.ProjectPath = rootDirectory
	c.localFlags.ArchFile = archFile

	return nil
}

func (c *CommandAssembler) assembleInput() (models.FlagsSelfInspect, error) {
	projectInfo, err := c.projectInfoAssembler.ProjectInfo(
		c.localFlags.ProjectPath,
		c.localFlags.ArchFile,
	)
	if err != nil {
		return models.FlagsSelfInspect{}, fmt.Errorf("failed to assemble project info: %w", err)
	}

	return models.FlagsSelfInspect{
		Project: projectInfo,
	}, nil
}
