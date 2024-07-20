package mapping

import (
	"github.com/urfave/cli/v2"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
	"github.com/fe3dback/go-arch-lint/v4/internal/pkg/xflags"
)

const (
	flagProjectPath            = "project-path"
	flagArchConfigRelativePath = "arch-file"
	flagScheme                 = "scheme"
)

var Flags = []cli.Flag{
	&cli.PathFlag{ // todo: add helper method with validation
		Name:     flagProjectPath,
		Category: models.FlagCategoryCommand,
		Usage:    "absolute path to project directory",
		Value:    models.DefaultProjectPath,
	},
	&cli.PathFlag{
		Name:     flagArchConfigRelativePath,
		Category: models.FlagCategoryCommand,
		Usage:    "relative path to linter config",
		Value:    models.DefaultArchFileName,
	},
	xflags.CreateEnumFlag(
		flagScheme,
		[]string{"s"},
		"display scheme",
		models.CmdMappingSchemesValues,
		models.CmdMappingSchemeList,
		models.FlagCategoryCommand,
	),
}
