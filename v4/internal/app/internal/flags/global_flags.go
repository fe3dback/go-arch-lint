package flags

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
	"github.com/fe3dback/go-arch-lint/v4/internal/pkg/xflags"
)

var GlobalFlags = []cli.Flag{
	xflags.CreateEnumFlag(
		models.FlagOutputType,
		[]string{},
		"linter output type",
		models.OutputTypeValues,
		models.OutputTypeASCII,
		models.FlagCategoryGlobal,
	),
	&cli.BoolFlag{
		Name:     models.FlagOutputTypeJSON,
		Category: models.FlagCategoryGlobal,
		Usage:    fmt.Sprintf("(alias for --%s %s)", models.FlagOutputType, models.OutputTypeJSON),
		Value:    false,
	},
	&cli.BoolFlag{
		Name:     models.FlagOutputJSONWithoutFormatting,
		Category: models.FlagCategoryGlobal,
		Usage:    fmt.Sprintf("output JSON in single line (without formatting). Only for '--%s'", models.FlagOutputTypeJSON),
		Value:    false,
	},
	&cli.BoolFlag{
		Name:     models.FlagOutputUseAsciiColors,
		Category: models.FlagCategoryGlobal,
		Usage:    "use ANSI colors in terminal output",
		Value:    true,
	},
	&cli.BoolFlag{
		Name:     models.FlagSkipMissUsages,
		Category: models.FlagCategoryGlobal,
		Usage:    "will skip not critical notices in config validation",
		Value:    false,
	},
	&cli.PathFlag{ // todo: add helper method with validation
		Name:     models.FlagProjectPath,
		Category: models.FlagCategoryCommand,
		Usage:    "absolute path to project directory",
		Value:    models.DefaultProjectPath,
	},
	&cli.PathFlag{
		Name:     models.FlagArchConfigRelativePath,
		Category: models.FlagCategoryCommand,
		Usage:    "relative path to linter config",
		Value:    models.DefaultArchFileName,
	},
}
