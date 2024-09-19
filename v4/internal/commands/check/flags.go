package check

import (
	"github.com/urfave/cli/v2"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

const (
	flagMaxWarnings = "max-warnings"
	flagSkipSyntax  = "skip-syntax"
)

var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name:     flagSkipSyntax,
		Category: models.FlagCategoryCommand,
		Usage:    "skip checking that golang code has valid syntax (other linters may not work as expected in this case)",
		Value:    false,
	},
	&cli.IntFlag{
		Name:     flagMaxWarnings,
		Category: models.FlagCategoryCommand,
		Usage:    "max number of warnings to output",
		Value:    100,
	},
}
