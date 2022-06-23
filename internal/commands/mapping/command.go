package mapping

import (
	"context"
	"fmt"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"

	"github.com/spf13/cobra"
)

const (
	flagProjectPath = "project-path"
	flagArchFile    = "arch-file"
	flagScheme      = "scheme"
)

var validSchemes = []string{
	models.MappingSchemeList,
	models.MappingSchemeGrouped,
}

type (
	processorFn = func(ctx context.Context, mapping models.FlagsMapping) error

	CommandAssembler struct {
		projectInfoAssembler ProjectInfoAssembler
		processorFn          processorFn
		localFlags           *localFlags
	}

	localFlags struct {
		ProjectPath string
		ArchFile    string
		Scheme      string
	}
)

func NewMappingCommandAssembler(
	projectInfoAssembler ProjectInfoAssembler,
	processorFn processorFn,
) *CommandAssembler {
	return &CommandAssembler{
		projectInfoAssembler: projectInfoAssembler,
		processorFn:          processorFn,
		localFlags: &localFlags{
			ProjectPath: "./",
			ArchFile:    models.ProjectInfoDefaultArchFileName,
			Scheme:      models.MappingSchemeList,
		},
	}
}

func (c *CommandAssembler) Assemble() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mapping",
		Aliases: []string{"ps", "ls"},
		Short:   "mapping table between files and components",
		Long:    "display mapping table between project files and arch components",
		PreRunE: c.prePersist,
		RunE:    c.invoke,
	}

	c.assembleFlags(cmd)

	return cmd
}

func (c *CommandAssembler) invoke(cmd *cobra.Command, _ []string) error {
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

	scheme, err := cmd.Flags().GetString(flagScheme)
	if err != nil {
		return failedToGetFlag(err, flagScheme)
	}

	schemeIsValid := false
	for _, s := range validSchemes {
		if scheme == s {
			schemeIsValid = true
			break
		}
	}

	if !schemeIsValid {
		return fmt.Errorf(
			"invalid scheme '%s', available: [%s]",
			scheme,
			strings.Join(validSchemes, ", "),
		)
	}

	// assemble localFlags
	c.localFlags.ProjectPath = rootDirectory
	c.localFlags.ArchFile = archFile
	c.localFlags.Scheme = scheme

	return nil
}

func (c *CommandAssembler) assembleInput() (models.FlagsMapping, error) {
	projectInfo, err := c.projectInfoAssembler.ProjectInfo(
		c.localFlags.ProjectPath,
		c.localFlags.ArchFile,
	)
	if err != nil {
		return models.FlagsMapping{}, fmt.Errorf("failed to assemble project info: %w", err)
	}

	return models.FlagsMapping{
		Project: projectInfo,
		Scheme:  c.localFlags.Scheme,
	}, nil
}
