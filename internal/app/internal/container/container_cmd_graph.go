package container

import (
	"fmt"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/operations/graph"
	"github.com/spf13/cobra"
)

func (c *Container) commandGraph() (*cobra.Command, runner) {
	cmd := &cobra.Command{
		Use:     "graph",
		Aliases: []string{"g"},
		Short:   "output dependencies graph as svg file",
		Long:    "display mapping table between project files and arch components",
	}

	in := models.CmdGraphIn{
		ProjectPath:    models.DefaultProjectPath,
		ArchFile:       models.DefaultArchFileName,
		Type:           models.GraphTypeFlow,
		OutFile:        "./go-arch-lint-graph.svg",
		Focus:          "",
		IncludeVendors: false,
		ExportD2:       false,
	}

	cmd.PersistentFlags().StringVar(&in.ProjectPath, "project-path", in.ProjectPath, "absolute path to project directory")
	cmd.PersistentFlags().StringVar(&in.ArchFile, "arch-file", in.ArchFile, "arch file path")
	cmd.PersistentFlags().StringVarP(&in.Type, "type", "t", in.Type, fmt.Sprintf("render graph type [%s]", strings.Join(models.GraphTypesValues, ",")))
	cmd.PersistentFlags().StringVar(&in.OutFile, "out", in.OutFile, "svg graph output file")
	cmd.PersistentFlags().StringVar(&in.Focus, "focus", in.Focus, "render only specified component (should match component name exactly)")
	cmd.PersistentFlags().BoolVarP(&in.IncludeVendors, "include-vendors", "r", in.IncludeVendors, "include vendor dependencies (from \"canUse\" block)?")
	cmd.PersistentFlags().BoolVar(&in.ExportD2, "d2", in.ExportD2, "output raw d2 definitions to stdout (from which svg is generated)")

	return cmd, func(act *cobra.Command) (any, error) {
		in.OutputType = c.flags.OutputType

		return c.commandGraphOperation().Behave(act.Context(), in)
	}
}

func (c *Container) commandGraphOperation() *graph.Operation {
	return graph.NewOperation(
		c.provideSpecAssembler(),
		c.provideProjectInfoAssembler(),
	)
}
