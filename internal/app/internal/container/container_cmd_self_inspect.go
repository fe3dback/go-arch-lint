package container

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/operations/selfInspect"
	"github.com/spf13/cobra"
)

func (c *Container) commandSelfInspect() (*cobra.Command, runner) {
	cmd := &cobra.Command{
		Use:   "self-inspect",
		Short: "will validate arch config and arch setup",
		Long:  "this useful for IDE plugins and other tool integration",
	}

	in := models.CmdSelfInspectIn{
		ProjectPath: models.DefaultProjectPath,
		ArchFile:    models.DefaultArchFileName,
	}

	cmd.PersistentFlags().StringVar(&in.ProjectPath, "project-path", in.ProjectPath, "absolute path to project directory")
	cmd.PersistentFlags().StringVar(&in.ArchFile, "arch-file", in.ArchFile, "arch file path")

	return cmd, func(_ *cobra.Command) (any, error) {
		return c.commandSelfInspectOperation().Behave(in)
	}
}

func (c *Container) commandSelfInspectOperation() *selfInspect.Operation {
	return selfInspect.NewOperation(
		c.provideSpecAssembler(),
		c.provideProjectInfoAssembler(),
		c.version,
	)
}
