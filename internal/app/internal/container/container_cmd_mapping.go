package container

import (
	"fmt"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/operations/mapping"
	"github.com/spf13/cobra"
)

func (c *Container) commandMapping() (*cobra.Command, runner) {
	cmd := &cobra.Command{
		Use:     "mapping",
		Aliases: []string{"ps", "ls"},
		Short:   "mapping table between files and components",
		Long:    "display mapping table between project files and arch components",
	}

	in := models.CmdMappingIn{
		ProjectPath: models.DefaultProjectPath,
		ArchFile:    models.DefaultArchFileName,
		Scheme:      models.MappingSchemeList,
	}

	cmd.PersistentFlags().StringVar(&in.ProjectPath, "project-path", in.ProjectPath, "absolute path to project directory")
	cmd.PersistentFlags().StringVar(&in.ArchFile, "arch-file", in.ArchFile, "arch file path")
	cmd.PersistentFlags().StringVarP(&in.Scheme, "scheme", "s", in.Scheme, fmt.Sprintf(
		"display scheme [%s]",
		strings.Join(models.MappingSchemesValues, ","),
	))

	return cmd, func(act *cobra.Command) (any, error) {
		hasValidScheme := false
		for _, validScheme := range models.MappingSchemesValues {
			if in.Scheme == validScheme {
				hasValidScheme = true
				break
			}
		}

		if !hasValidScheme {
			return "", fmt.Errorf(
				"invalid scheme '%s', available: [%s]",
				in.Scheme,
				strings.Join(models.MappingSchemesValues, ", "),
			)
		}

		return c.commandMappingOperation().Behave(act.Context(), in)
	}
}

func (c *Container) commandMappingOperation() *mapping.Operation {
	return mapping.NewOperation(
		c.provideSpecAssembler(),
		c.provideProjectFilesResolver(),
		c.provideProjectInfoAssembler(),
	)
}
