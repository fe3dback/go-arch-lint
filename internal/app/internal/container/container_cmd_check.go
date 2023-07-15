package container

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/operations/check"
	"github.com/spf13/cobra"
)

func (c *Container) commandCheck() (*cobra.Command, runner) {
	cmd := &cobra.Command{
		Use:     "check",
		Aliases: []string{"c"},
		Short:   "check project architecture by yaml file",
		Long:    "compare project *.go files with arch defined in spec file",
	}

	in := models.CmdCheckIn{
		ProjectPath: models.DefaultProjectPath,
		ArchFile:    models.DefaultArchFileName,
		MaxWarnings: 512,
	}

	cmd.PersistentFlags().StringVar(&in.ProjectPath, "project-path", in.ProjectPath, "absolute path to project directory")
	cmd.PersistentFlags().StringVar(&in.ArchFile, "arch-file", in.ArchFile, "arch file path")
	cmd.PersistentFlags().IntVar(&in.MaxWarnings, "max-warnings", in.MaxWarnings, "max number of warnings to output")

	return cmd, func(act *cobra.Command) (any, error) {
		const warningsRangeMin = 1
		const warningsRangeMax = 32768

		if in.MaxWarnings < warningsRangeMin || in.MaxWarnings > warningsRangeMax {
			return nil, fmt.Errorf(
				"flag '%s' should by in range [%d .. %d]",
				"max-warnings",
				warningsRangeMin,
				warningsRangeMax,
			)
		}

		return c.commandCheckOperation().Behave(act.Context(), in)
	}
}

func (c *Container) commandCheckOperation() *check.Operation {
	return check.NewOperation(
		c.provideProjectInfoAssembler(),
		c.provideSpecAssembler(),
		c.provideSpecChecker(),
		c.provideReferenceRender(),
		c.flags.UseColors,
	)
}
