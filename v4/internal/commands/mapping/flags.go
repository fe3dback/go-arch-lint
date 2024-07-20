package mapping

import (
	"github.com/urfave/cli/v2"

	"github.com/fe3dback/go-arch-lint/v4/internal/commands"
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

const (
	flagProjectPath            = "project-path"
	flagArchConfigRelativePath = "arch-file"
	flagScheme                 = "scheme"
)

var Flags = []cli.Flag{
	&cli.PathFlag{ // todo: add helper method with validation
		Name:  flagProjectPath,
		Usage: "absolute path to project directory",
		Value: commands.DefaultProjectPath,
	},
	&cli.PathFlag{
		Name:  flagArchConfigRelativePath,
		Usage: "relative path to linter config",
		Value: commands.DefaultArchFileName,
	},
	commands.CreateEnumFlag(
		flagScheme,
		[]string{"s"},
		"display scheme",
		models.CmdMappingSchemesValues,
		models.CmdMappingSchemeList,
	),
}
