package check

import (
	"github.com/urfave/cli/v2"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

const (
	flagMaxWarnings = "max-warnings"
)

var Flags = []cli.Flag{
	&cli.IntFlag{
		Name:     flagMaxWarnings,
		Category: models.FlagCategoryCommand,
		Usage:    "max number of warnings to output",
		Value:    100,
	},
}
