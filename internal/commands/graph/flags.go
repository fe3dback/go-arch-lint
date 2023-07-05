package graph

import (
	"fmt"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"

	"github.com/spf13/cobra"
)

const (
	flagProjectPath   = "project-path"
	flagArchFile      = "arch-file"
	flagOutFile       = "out"
	flagGraphType     = "type"
	flagFocus         = "focus"
	flagIncludeVendor = "include-vendors"
)

type (
	localFlags struct {
		ProjectPath   string
		ArchFile      string
		OutFile       string
		GraphType     string
		Focus         string
		IncludeVendor bool
	}
)

var validGraphTypes = []string{
	models.GraphTypeFlow,
	models.GraphTypeDI,
}

func (c *CommandAssembler) assembleFlags(cmd *cobra.Command) {
	c.withProjectPath(cmd)
	c.withArchFileName(cmd)
	c.withOutFileName(cmd)
	c.withGraphType(cmd)
	c.withFocus(cmd)
	c.withIncludeVendors(cmd)
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

func (c *CommandAssembler) withOutFileName(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		flagOutFile,
		"",
		c.localFlags.OutFile,
		"svg graph output file",
	)
}

func (c *CommandAssembler) withGraphType(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		flagGraphType,
		"t",
		c.localFlags.GraphType,
		fmt.Sprintf(
			"render graph type [%s]",
			strings.Join(validGraphTypes, ","),
		),
	)
}

func (c *CommandAssembler) withFocus(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		flagFocus,
		"",
		c.localFlags.Focus,
		"render only specified component (should match component name exactly)",
	)
}

func (c *CommandAssembler) withIncludeVendors(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(
		flagIncludeVendor,
		"r",
		c.localFlags.IncludeVendor,
		"include vendor dependencies (from \"canUse\" block)?",
	)
}
