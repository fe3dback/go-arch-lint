package graph

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"

	"github.com/spf13/cobra"
)

type (
	processorFn = func(ctx context.Context, graph models.FlagsGraph) error

	CommandAssembler struct {
		projectInfoAssembler ProjectInfoAssembler
		processorFn          processorFn
		localFlags           *localFlags
	}
)

func NewGraphCommandAssembler(
	projectInfoAssembler ProjectInfoAssembler,
	processorFn processorFn,
) *CommandAssembler {
	return &CommandAssembler{
		projectInfoAssembler: projectInfoAssembler,
		processorFn:          processorFn,
		localFlags: &localFlags{
			ProjectPath: "./",
			ArchFile:    models.ProjectInfoDefaultArchFileName,
			OutFile:     "./go-arch-lint-graph.svg",
			GraphType:   models.GraphTypeFlow,
		},
	}
}

func (c *CommandAssembler) Assemble() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "graph",
		Aliases: []string{"g"},
		Short:   "output dependencies graph as svg file",
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

	outFileRel, err := cmd.Flags().GetString(flagOutFile)
	if err != nil {
		return failedToGetFlag(err, flagOutFile)
	}

	outFile, err := filepath.Abs(outFileRel)
	if err != nil {
		return fmt.Errorf("failed get abs path from '%s': %w", outFileRel, err)
	}

	err = os.WriteFile(outFile, nil, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write 'out' '%s' file into '%s': %w", outFileRel, outFile, err)
	}

	graphType, err := cmd.Flags().GetString(flagGraphType)
	if err != nil {
		return failedToGetFlag(err, flagGraphType)
	}

	graphTypeIsValid := false
	for _, s := range validGraphTypes {
		if graphType == s {
			graphTypeIsValid = true
			break
		}
	}

	if !graphTypeIsValid {
		return fmt.Errorf(
			"invalid graph type '%s', available: [%s]",
			graphType,
			strings.Join(validGraphTypes, ", "),
		)
	}

	includeVendors, err := cmd.Flags().GetBool(flagIncludeVendor)
	if err != nil {
		return failedToGetFlag(err, flagIncludeVendor)
	}

	focus, err := cmd.Flags().GetString(flagFocus)
	if err != nil {
		return failedToGetFlag(err, flagFocus)
	}

	// assemble localFlags
	c.localFlags.ProjectPath = rootDirectory
	c.localFlags.ArchFile = archFile
	c.localFlags.OutFile = outFile
	c.localFlags.GraphType = graphType
	c.localFlags.IncludeVendor = includeVendors
	c.localFlags.Focus = focus

	return nil
}

func (c *CommandAssembler) assembleInput() (models.FlagsGraph, error) {
	projectInfo, err := c.projectInfoAssembler.ProjectInfo(
		c.localFlags.ProjectPath,
		c.localFlags.ArchFile,
	)
	if err != nil {
		return models.FlagsGraph{}, fmt.Errorf("failed to assemble project info: %w", err)
	}

	return models.FlagsGraph{
		Project:        projectInfo,
		OutFile:        c.localFlags.OutFile,
		Type:           c.localFlags.GraphType,
		IncludeVendors: c.localFlags.IncludeVendor,
		Focus:          c.localFlags.Focus,
	}, nil
}
